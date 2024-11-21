use super::super::operator::wrap;
use crate::compactor::CompactionJob;
use crate::execution::dispatcher::Dispatcher;
use crate::execution::operator::TaskResult;
use crate::execution::operators::apply_log_to_segment_writer::{
    ApplyLogToSegmentWriterInput, ApplyLogToSegmentWriterOperator,
    ApplyLogToSegmentWriterOperatorError, ApplyLogToSegmentWriterOutput,
};
use crate::execution::operators::commit_segment_writer::{
    CommitSegmentWriterInput, CommitSegmentWriterOperator, CommitSegmentWriterOperatorError,
    CommitSegmentWriterOutput,
};
use crate::execution::operators::flush_segment_writer::{
    FlushSegmentWriterInput, FlushSegmentWriterOperator, FlushSegmentWriterOperatorError,
    FlushSegmentWriterOutput,
};
use crate::execution::operators::materialize_logs::MaterializeLogInput;
use crate::execution::operators::materialize_logs::MaterializeLogOperator;
use crate::execution::operators::materialize_logs::MaterializeLogOperatorError;
use crate::execution::operators::materialize_logs::MaterializeLogOutput;
use crate::execution::operators::partition::PartitionError;
use crate::execution::operators::partition::PartitionInput;
use crate::execution::operators::partition::PartitionOperator;
use crate::execution::operators::partition::PartitionOutput;
use crate::execution::operators::pull_log::PullLogsInput;
use crate::execution::operators::pull_log::PullLogsOperator;
use crate::execution::operators::pull_log::PullLogsOutput;
use crate::execution::operators::register::RegisterError;
use crate::execution::operators::register::RegisterInput;
use crate::execution::operators::register::RegisterOperator;
use crate::execution::operators::register::RegisterOutput;
use crate::execution::orchestration::common::terminate_with_error;
use crate::log::log::Log;
use crate::log::log::PullLogsError;
use crate::segment::distributed_hnsw_segment::DistributedHNSWSegmentWriter;
use crate::segment::metadata_segment::MetadataSegmentWriter;
use crate::segment::record_segment::RecordSegmentReader;
use crate::segment::record_segment::RecordSegmentReaderCreationError;
use crate::segment::record_segment::RecordSegmentWriter;
use crate::segment::{ChromaSegmentFlusher, ChromaSegmentWriter, SegmentFlusher, SegmentWriter};
use crate::sysdb::sysdb::GetCollectionsError;
use crate::sysdb::sysdb::GetSegmentsError;
use crate::sysdb::sysdb::SysDb;
use crate::system::Component;
use crate::system::ComponentHandle;
use crate::system::Handler;
use crate::system::ReceiverForMessage;
use crate::system::System;
use async_trait::async_trait;
use chroma_blockstore::provider::BlockfileProvider;
use chroma_error::ChromaError;
use chroma_error::ErrorCodes;
use chroma_index::hnsw_provider::HnswIndexProvider;
use chroma_types::{Chunk, SegmentUuid};
use chroma_types::{CollectionUuid, LogRecord, Segment, SegmentFlushInfo, SegmentType};
use core::panic;
use std::collections::HashMap;
use std::sync::atomic::AtomicU32;
use std::sync::Arc;
use std::time::SystemTime;
use std::time::UNIX_EPOCH;
use thiserror::Error;
use tokio::sync::OnceCell;
use tracing::Span;
use uuid::Uuid;

/**  The state of the orchestrator.
In chroma, we have a relatively fixed number of query plans that we can execute. Rather
than a flexible state machine abstraction, we just manually define the states that we
expect to encounter for a given query plan. This is a bit more rigid, but it's also simpler and easier to
understand. We can always add more abstraction later if we need it.
```plaintext
                                   ┌────────────────────────────┐
                                   ├─► Apply logs to segment #1 │
                                   │                            ├──► Commit segment #1 ──► Flush segment #1
                                   ├─► Apply logs to segment #1 │
Pending ──► PullLogs ──► Partition │                            │                                            ──► Register ─► Finished
                                   ├─► Apply logs to segment #2 │
                                   │                            ├──► Commit segment #2 ──► Flush segment #2
                                   ├─► Apply logs to segment #2 │
                                   └────────────────────────────┘
```
*/
#[derive(Debug)]
enum ExecutionState {
    Pending,
    PullLogs,
    Partition,
    MaterializeApplyCommitFlush,
    Register,
}

#[derive(Debug)]
pub struct CompactOrchestrator {
    id: Uuid,
    compaction_job: CompactionJob,
    state: ExecutionState,
    // Component Execution
    system: System,
    collection_id: CollectionUuid,
    // Dependencies
    log: Box<Log>,
    sysdb: Box<SysDb>,
    blockfile_provider: BlockfileProvider,
    hnsw_index_provider: HnswIndexProvider,
    // State we hold across the execution
    pulled_log_offset: Option<i64>,
    // Dispatcher
    dispatcher: ComponentHandle<Dispatcher>,
    // Tracks the total remaining number of ApplyLogToSegmentWriter tasks per segment
    num_uncompleted_log_apply_tasks: HashMap<SegmentUuid, usize>,
    // Tracks the total remaining number of FlushSegmentWriter tasks
    num_uncompleted_flush_tasks: usize,
    // Result Channel
    result_channel:
        Option<tokio::sync::oneshot::Sender<Result<CompactionResponse, Box<dyn ChromaError>>>>,
    max_compaction_size: usize,
    max_partition_size: usize,
    // Populated during the compaction process
    cached_segments: OnceCell<Vec<Segment>>,
    writers: OnceCell<(
        RecordSegmentWriter,
        Box<DistributedHNSWSegmentWriter>,
        MetadataSegmentWriter<'static>,
    )>,
    flush_results: Vec<SegmentFlushInfo>,
    // We track a parent span for each segment type so we can group all the spans for a given segment type (makes the resulting trace much easier to read)
    segment_spans: HashMap<SegmentUuid, Span>,
}

#[derive(Error, Debug)]
enum GetSegmentWritersError {
    #[error("No segments found for collection")]
    NoSegmentsFound,
    #[error("SysDB GetSegments Error")]
    SysDbGetSegmentsError(#[from] GetSegmentsError),
    #[error("Error creating Record Segment Writer")]
    RecordSegmentWriterError,
    #[error("Error creating Metadata Segment Writer")]
    MetadataSegmentWriterError,
    #[error("Error creating HNSW Segment Writer")]
    HnswSegmentWriterError,
    #[error("Collection not found")]
    CollectionNotFound,
    #[error("Error getting collection")]
    GetCollectionError(#[from] GetCollectionsError),
}

impl ChromaError for GetSegmentWritersError {
    fn code(&self) -> ErrorCodes {
        ErrorCodes::Internal
    }
}

#[derive(Error, Debug)]
enum CompactionError {
    #[error(transparent)]
    SystemTimeError(#[from] std::time::SystemTimeError),
    #[error("Result channel dropped")]
    ResultChannelDropped,
    #[error("Invariant violation: {}", .0)]
    InvariantViolation(&'static str),
}

impl ChromaError for CompactionError {
    fn code(&self) -> ErrorCodes {
        ErrorCodes::Internal
    }
}

// TODO: we need to improve this response
#[derive(Debug)]
pub struct CompactionResponse {
    #[allow(dead_code)]
    pub(crate) id: Uuid,
    pub(crate) compaction_job: CompactionJob,
    #[allow(dead_code)]
    pub(crate) message: String,
}

impl CompactOrchestrator {
    #[allow(clippy::too_many_arguments)]
    pub fn new(
        compaction_job: CompactionJob,
        system: System,
        collection_id: CollectionUuid,
        log: Box<Log>,
        sysdb: Box<SysDb>,
        blockfile_provider: BlockfileProvider,
        hnsw_index_provider: HnswIndexProvider,
        dispatcher: ComponentHandle<Dispatcher>,
        result_channel: Option<
            tokio::sync::oneshot::Sender<Result<CompactionResponse, Box<dyn ChromaError>>>,
        >,
        max_compaction_size: usize,
        max_partition_size: usize,
    ) -> Self {
        CompactOrchestrator {
            id: Uuid::new_v4(),
            compaction_job,
            state: ExecutionState::Pending,
            system,
            collection_id,
            log,
            sysdb,
            blockfile_provider,
            hnsw_index_provider,
            pulled_log_offset: None,
            dispatcher,
            num_uncompleted_log_apply_tasks: HashMap::new(),
            num_uncompleted_flush_tasks: 0,
            result_channel,
            max_compaction_size,
            max_partition_size,
            cached_segments: OnceCell::new(),
            writers: OnceCell::new(),
            flush_results: Vec::new(),
            segment_spans: HashMap::new(),
        }
    }

    async fn pull_logs(
        &mut self,
        self_address: Box<dyn ReceiverForMessage<TaskResult<PullLogsOutput, PullLogsError>>>,
        ctx: &crate::system::ComponentContext<CompactOrchestrator>,
    ) {
        self.state = ExecutionState::PullLogs;
        let operator = PullLogsOperator::new(self.log.clone());
        let collection_id = self.collection_id;
        let end_timestamp = SystemTime::now().duration_since(UNIX_EPOCH);
        let end_timestamp = match end_timestamp {
            // TODO: change protobuf definition to use u64 instead of i64
            Ok(end_timestamp) => end_timestamp.as_nanos() as i64,
            Err(e) => {
                terminate_with_error(
                    self.result_channel.take(),
                    Box::new(CompactionError::SystemTimeError(e)),
                    ctx,
                );
                return;
            }
        };
        let input = PullLogsInput::new(
            collection_id,
            // Here we do not need to be inclusive since the compaction job
            // offset is the one after the last compaction offset
            self.compaction_job.offset,
            100,
            Some(self.max_compaction_size as i32),
            Some(end_timestamp),
        );
        let task = wrap(operator, input, self_address);
        match self.dispatcher.send(task, Some(Span::current())).await {
            Ok(_) => (),
            Err(e) => {
                tracing::error!("Error dispatching pull logs for compaction {:?}", e);
                panic!(
                    "Invariant violation. Somehow the dispatcher receiver is dropped. Error: {:?}",
                    e
                );
            }
        }
    }

    async fn partition(
        &mut self,
        records: Chunk<LogRecord>,
        self_address: Box<dyn ReceiverForMessage<TaskResult<PartitionOutput, PartitionError>>>,
    ) {
        self.state = ExecutionState::Partition;
        let operator = PartitionOperator::new();
        tracing::info!("Sending N Records: {:?}", records.len());
        println!("Sending N Records: {:?}", records.len());
        let input = PartitionInput::new(records, self.max_partition_size);
        let task = wrap(operator, input, self_address);
        match self.dispatcher.send(task, Some(Span::current())).await {
            Ok(_) => (),
            Err(e) => {
                tracing::error!("Error dispatching partition for compaction {:?}", e);
                panic!(
                    "Invariant violation. Somehow the dispatcher receiver is dropped. Error: {:?}",
                    e
                )
            }
        }
    }

    async fn materialize_log(
        &mut self,
        partitions: Vec<Chunk<LogRecord>>,
        self_address: Box<
            dyn ReceiverForMessage<TaskResult<MaterializeLogOutput, MaterializeLogOperatorError>>,
        >,
        ctx: &crate::system::ComponentContext<CompactOrchestrator>,
    ) {
        self.state = ExecutionState::MaterializeApplyCommitFlush;

        let record_segment = match self.get_segment(SegmentType::BlockfileRecord).await {
            Ok(segment) => segment,
            Err(e) => {
                tracing::error!("Error getting record segment: {:?}", e);
                terminate_with_error(self.result_channel.take(), e, ctx);
                return;
            }
        };

        let hnsw_segment = match self.get_segment(SegmentType::HnswDistributed).await {
            Ok(segment) => segment,
            Err(e) => {
                tracing::error!("Error getting hnsw segment: {:?}", e);
                terminate_with_error(self.result_channel.take(), e, ctx);
                return;
            }
        };

        let metadata_segment = match self.get_segment(SegmentType::BlockfileMetadata).await {
            Ok(segment) => segment,
            Err(e) => {
                tracing::error!("Error getting metadata segment: {:?}", e);
                terminate_with_error(self.result_channel.take(), e, ctx);
                return;
            }
        };

        let current_max_offset_id = match RecordSegmentReader::from_segment(
            &record_segment,
            &self.blockfile_provider,
        )
        .await
        {
            Ok(reader) => Arc::new(AtomicU32::new(reader.get_max_offset_id())),
            Err(err) => match *err {
                RecordSegmentReaderCreationError::UninitializedSegment => {
                    Arc::new(AtomicU32::new(0))
                }
                _ => {
                    tracing::error!("Error getting current max offset id: {:?}", err);
                    terminate_with_error(
                        self.result_channel.take(),
                        err as Box<dyn ChromaError>,
                        ctx,
                    );
                    return;
                }
            },
        };

        self.num_uncompleted_log_apply_tasks
            .insert(record_segment.id, partitions.len());
        self.num_uncompleted_log_apply_tasks
            .insert(hnsw_segment.id, partitions.len());
        self.num_uncompleted_log_apply_tasks
            .insert(metadata_segment.id, partitions.len());
        self.num_uncompleted_flush_tasks = 3;

        for partition in partitions.iter() {
            let operator = MaterializeLogOperator::new();
            let input = MaterializeLogInput::new(
                partition.clone(),
                self.blockfile_provider.clone(),
                record_segment.clone(),
                current_max_offset_id.clone(),
            );
            let task = wrap(operator, input, self_address.clone());
            match self.dispatcher.send(task, Some(Span::current())).await {
                Ok(_) => (),
                Err(e) => {
                    tracing::error!(
                        "Error dispatching log materialization tasks for compaction {:?}",
                        e
                    );
                    panic!(
                        "Invariant violation. Somehow the dispatcher receiver is dropped. Error: {:?}",
                        e)
                }
            }
        }
    }

    fn get_segment_writer_span(&mut self, writer: &ChromaSegmentWriter<'static>) -> Span {
        let span = self
            .segment_spans
            .entry(writer.get_id())
            .or_insert_with(|| {
                tracing::span!(
                    tracing::Level::INFO,
                    "Segment",
                    otel.name = format!("Segment: {:?}", writer.get_name())
                )
            });
        span.clone()
    }

    fn get_segment_flusher_span(&mut self, flusher: &ChromaSegmentFlusher) -> Span {
        match self.segment_spans.get(&flusher.get_id()) {
            Some(span) => span.clone(),
            None => {
                tracing::error!(
                    "No span found for segment: {:?}. This should never happen because get_segment_writer_span() should have previously created a span.",
                    flusher.get_name()
                );
                Span::current()
            }
        }
    }

    async fn dispatch_apply_log_to_segment_writer_task(
        &mut self,
        segment_writer: ChromaSegmentWriter<'static>,
        materialized_log: Arc<MaterializeLogOutput>,
        self_address: Box<
            dyn ReceiverForMessage<
                TaskResult<
                    ApplyLogToSegmentWriterOutput<'static>,
                    ApplyLogToSegmentWriterOperatorError,
                >,
            >,
        >,
    ) {
        let span = self.get_segment_writer_span(&segment_writer);
        let operator = ApplyLogToSegmentWriterOperator::new();
        let input = ApplyLogToSegmentWriterInput::new(segment_writer, materialized_log);
        let task = wrap(operator, input, self_address.clone());
        match self.dispatcher.send(task, Some(span)).await {
            Ok(_) => {}
            Err(e) => {
                tracing::error!(
                    "Error dispatching apply log to segment writer task: {:?}",
                    e
                );
                panic!(
                    "Invariant violation. Somehow the dispatcher receiver is dropped. Error: {:?}",
                    e
                )
            }
        }
    }

    async fn dispatch_segment_writer_commit(
        &mut self,
        segment_writer: ChromaSegmentWriter<'static>,
        self_address: Box<
            dyn ReceiverForMessage<
                TaskResult<CommitSegmentWriterOutput, CommitSegmentWriterOperatorError>,
            >,
        >,
    ) {
        let span = self.get_segment_writer_span(&segment_writer);
        let operator = CommitSegmentWriterOperator::new();
        let input = CommitSegmentWriterInput::new(segment_writer);
        let task = wrap(operator, input, self_address);
        match self.dispatcher.send(task, Some(span)).await {
            Ok(_) => (),
            Err(e) => {
                tracing::error!("Error dispatching commit for compaction {:?}", e);
                panic!(
                    "Invariant violation. Somehow the dispatcher receiver is dropped. Error: {:?}",
                    e
                );
            }
        }
    }

    async fn dispatch_segment_flush(
        &mut self,
        segment_flusher: ChromaSegmentFlusher,
        self_address: Box<
            dyn ReceiverForMessage<
                TaskResult<FlushSegmentWriterOutput, FlushSegmentWriterOperatorError>,
            >,
        >,
    ) {
        let span = self.get_segment_flusher_span(&segment_flusher);
        let operator = FlushSegmentWriterOperator::new();
        let input = FlushSegmentWriterInput::new(segment_flusher);
        let task = wrap(operator, input, self_address);
        match self.dispatcher.send(task, Some(span)).await {
            Ok(_) => {}
            Err(e) => {
                tracing::error!("Error dispatching flush for compaction {:?}", e);
                panic!(
                    "Invariant violation. Somehow the dispatcher receiver is dropped. Error: {:?}",
                    e
                );
            }
        }
    }

    async fn register(
        &mut self,
        log_position: i64,
        self_address: Box<dyn ReceiverForMessage<TaskResult<RegisterOutput, RegisterError>>>,
    ) {
        self.state = ExecutionState::Register;
        let operator = RegisterOperator::new();
        let input = RegisterInput::new(
            self.compaction_job.tenant_id.clone(),
            self.compaction_job.collection_id,
            log_position,
            self.compaction_job.collection_version,
            self.flush_results.clone().into(),
            self.sysdb.clone(),
            self.log.clone(),
        );

        let task = wrap(operator, input, self_address);
        match self.dispatcher.send(task, Some(Span::current())).await {
            Ok(_) => (),
            Err(e) => {
                tracing::error!("Error dispatching register for compaction {:?}", e);
                panic!(
                    "Invariant violation. Somehow the dispatcher receiver is dropped. Error: {:?}",
                    e
                );
            }
        }
    }

    async fn get_all_segments(&mut self) -> Result<Vec<Segment>, GetSegmentsError> {
        self.cached_segments
            .get_or_try_init(|| async {
                self.sysdb
                    .get_segments(None, None, None, self.collection_id)
                    .await
            })
            .await
            .cloned()
    }

    async fn get_segment(
        &mut self,
        segment_type: SegmentType,
    ) -> Result<Segment, Box<dyn ChromaError>> {
        let segments = self
            .get_all_segments()
            .await
            .map_err(|e| Box::new(e) as Box<dyn ChromaError>)?;
        let segment = segments
            .iter()
            .find(|segment| segment.r#type == segment_type)
            .cloned();

        tracing::debug!("Found {:?} segment: {:?}", segment_type, segment);

        match segment {
            Some(segment) => Ok(segment),
            None => Err(Box::new(GetSegmentWritersError::NoSegmentsFound)),
        }
    }

    async fn get_segment_writers(
        &mut self,
    ) -> Result<
        (
            RecordSegmentWriter,
            Box<DistributedHNSWSegmentWriter>,
            MetadataSegmentWriter<'static>,
        ),
        Box<dyn ChromaError>,
    > {
        // Care should be taken to use the same writers across the compaction process
        // Since the segment writers are stateful, we should not create new writers for each partition
        // Nor should we create new writers across different tasks

        let blockfile_provider = self.blockfile_provider.clone();
        let hnsw_provider = self.hnsw_index_provider.clone();
        let mut sysdb = self.sysdb.clone();

        let record_segment = self.get_segment(SegmentType::BlockfileRecord).await?;
        let mt_segment = self.get_segment(SegmentType::BlockfileMetadata).await?;
        let hnsw_segment = self.get_segment(SegmentType::HnswDistributed).await?;

        let borrowed_writers = self
            .writers
            .get_or_try_init::<Box<dyn ChromaError>, _, _>(|| async {
                // Create a record segment writer
                let record_segment_writer =
                    match RecordSegmentWriter::from_segment(&record_segment, &blockfile_provider)
                        .await
                    {
                        Ok(writer) => writer,
                        Err(e) => {
                            tracing::error!("Error creating Record Segment Writer: {:?}", e);
                            return Err(Box::new(GetSegmentWritersError::RecordSegmentWriterError)
                                as Box<dyn ChromaError>);
                        }
                    };

                tracing::debug!("Record Segment Writer created");

                // Create a record segment writer
                let mt_segment_writer =
                    match MetadataSegmentWriter::from_segment(&mt_segment, &blockfile_provider)
                        .await
                    {
                        Ok(writer) => writer,
                        Err(e) => {
                            println!("Error creating metadata Segment Writer: {:?}", e);
                            return Err(Box::new(
                                GetSegmentWritersError::MetadataSegmentWriterError,
                            ));
                        }
                    };

                tracing::debug!("Metadata Segment Writer created");

                // Create a hnsw segment writer
                let collection_res = sysdb
                    .get_collections(Some(self.collection_id), None, None, None)
                    .await;

                let collection_res = match collection_res {
                    Ok(collections) => {
                        if collections.is_empty() {
                            return Err(Box::new(GetSegmentWritersError::CollectionNotFound));
                        }
                        collections
                    }
                    Err(e) => {
                        return Err(Box::new(GetSegmentWritersError::GetCollectionError(e)));
                    }
                };
                let collection = &collection_res[0];

                let dimension = collection
                    .dimension
                    .expect("Dimension is required in the compactor");

                let hnsw_segment_writer = match DistributedHNSWSegmentWriter::from_segment(
                    &hnsw_segment,
                    dimension as usize,
                    hnsw_provider,
                )
                .await
                {
                    Ok(writer) => writer,
                    Err(e) => {
                        println!("Error creating HNSW Segment Writer: {:?}", e);
                        return Err(Box::new(GetSegmentWritersError::HnswSegmentWriterError));
                    }
                };

                Ok((
                    record_segment_writer,
                    hnsw_segment_writer,
                    mt_segment_writer,
                ))
            })
            .await?;

        Ok(borrowed_writers.clone())
    }

    pub(crate) async fn run(mut self) -> Result<CompactionResponse, Box<dyn ChromaError>> {
        println!("Running compaction job: {:?}", self.compaction_job);
        let (tx, rx) = tokio::sync::oneshot::channel();
        self.result_channel = Some(tx);
        let mut handle = self.system.clone().start_component(self);
        let result = rx.await;
        handle.stop();
        result
            .map_err(|_| Box::new(CompactionError::ResultChannelDropped) as Box<dyn ChromaError>)?
    }
}

// ============== Component Implementation ==============

#[async_trait]
impl Component for CompactOrchestrator {
    fn get_name() -> &'static str {
        "Compaction orchestrator"
    }

    fn queue_size(&self) -> usize {
        1000 // TODO: make configurable
    }

    async fn on_start(&mut self, ctx: &crate::system::ComponentContext<Self>) -> () {
        self.pull_logs(ctx.receiver(), ctx).await;
    }
}

// ============== Handlers ==============
#[async_trait]
impl Handler<TaskResult<PullLogsOutput, PullLogsError>> for CompactOrchestrator {
    type Result = ();

    async fn handle(
        &mut self,
        message: TaskResult<PullLogsOutput, PullLogsError>,
        ctx: &crate::system::ComponentContext<CompactOrchestrator>,
    ) {
        let message = message.into_inner();
        let records = match message {
            Ok(result) => result.logs(),
            Err(e) => {
                terminate_with_error(self.result_channel.take(), Box::new(e), ctx);
                return;
            }
        };
        tracing::info!("Pulled Records: {:?}", records.len());
        let final_record_pulled = records.get(records.len() - 1);
        match final_record_pulled {
            Some(record) => {
                self.pulled_log_offset = Some(record.log_offset);
                tracing::info!("Pulled Logs Up To Offset: {:?}", self.pulled_log_offset);
                self.partition(records, ctx.receiver()).await;
            }
            None => {
                tracing::error!(
                    "No records pulled by compaction, this is a system invariant violation"
                );
                panic!("No records pulled by compaction, this is a system invariant violation");
            }
        }
    }
}

#[async_trait]
impl Handler<TaskResult<PartitionOutput, PartitionError>> for CompactOrchestrator {
    type Result = ();

    async fn handle(
        &mut self,
        message: TaskResult<PartitionOutput, PartitionError>,
        ctx: &crate::system::ComponentContext<CompactOrchestrator>,
    ) {
        let message = message.into_inner();
        let records = match message {
            Ok(result) => result.records,
            Err(e) => {
                tracing::error!("Error partitioning records: {:?}", e);
                terminate_with_error(self.result_channel.take(), Box::new(e), ctx);
                return;
            }
        };
        self.materialize_log(records, ctx.receiver(), ctx).await;
    }
}

#[async_trait]
impl Handler<TaskResult<MaterializeLogOutput, MaterializeLogOperatorError>>
    for CompactOrchestrator
{
    type Result = ();

    async fn handle(
        &mut self,
        message: TaskResult<MaterializeLogOutput, MaterializeLogOperatorError>,
        ctx: &crate::system::ComponentContext<CompactOrchestrator>,
    ) {
        let message = message.into_inner();
        let materialized_result = match message {
            Ok(result) => result,
            Err(e) => {
                tracing::error!("Error materializing log: {:?}", e);
                terminate_with_error(self.result_channel.take(), Box::new(e), ctx);
                return;
            }
        };

        let materialized_result = Arc::new(materialized_result);

        let (record_segment_writer, hnsw_segment_writer, metadata_segment_writer) =
            match self.get_segment_writers().await {
                Ok(writers) => writers,
                Err(e) => {
                    tracing::error!("Error getting segment writers: {:?}", e);
                    terminate_with_error(self.result_channel.take(), e, ctx);
                    return;
                }
            };

        self.dispatch_apply_log_to_segment_writer_task(
            ChromaSegmentWriter::RecordSegment(record_segment_writer),
            materialized_result.clone(),
            ctx.receiver(),
        )
        .await;

        self.dispatch_apply_log_to_segment_writer_task(
            ChromaSegmentWriter::DistributedHNSWSegment(hnsw_segment_writer),
            materialized_result.clone(),
            ctx.receiver(),
        )
        .await;

        self.dispatch_apply_log_to_segment_writer_task(
            ChromaSegmentWriter::MetadataSegment(metadata_segment_writer),
            materialized_result,
            ctx.receiver(),
        )
        .await;
    }
}

#[async_trait]
impl
    Handler<
        TaskResult<ApplyLogToSegmentWriterOutput<'static>, ApplyLogToSegmentWriterOperatorError>,
    > for CompactOrchestrator
{
    type Result = ();

    async fn handle(
        &mut self,
        message: TaskResult<
            ApplyLogToSegmentWriterOutput<'static>,
            ApplyLogToSegmentWriterOperatorError,
        >,
        ctx: &crate::system::ComponentContext<CompactOrchestrator>,
    ) {
        let message = message.into_inner();
        match message {
            Ok(message) => {
                self.num_uncompleted_log_apply_tasks
                    .entry(message.segment_writer.get_id())
                    .and_modify(|e| *e -= 1);

                let num_tasks_left = match self
                    .num_uncompleted_log_apply_tasks
                    .get(&message.segment_writer.get_id())
                {
                    Some(num_tasks_left) => num_tasks_left,
                    None => {
                        terminate_with_error(
                            self.result_channel.take(),
                            Box::new(CompactionError::InvariantViolation(
                                "Segment writer not found",
                            )),
                            ctx,
                        );
                        return;
                    }
                };

                if *num_tasks_left == 0 {
                    self.dispatch_segment_writer_commit(message.segment_writer, ctx.receiver())
                        .await;
                }
            }
            Err(e) => {
                tracing::error!("Error writing segments: {:?}", e);
                terminate_with_error(self.result_channel.take(), Box::new(e), ctx);
            }
        }
    }
}

#[async_trait]
impl Handler<TaskResult<CommitSegmentWriterOutput, CommitSegmentWriterOperatorError>>
    for CompactOrchestrator
{
    type Result = ();

    async fn handle(
        &mut self,
        message: TaskResult<CommitSegmentWriterOutput, CommitSegmentWriterOperatorError>,
        ctx: &crate::system::ComponentContext<CompactOrchestrator>,
    ) {
        let message = message.into_inner();
        match message {
            Ok(message) => {
                self.dispatch_segment_flush(message.flusher, ctx.receiver())
                    .await;
            }
            Err(e) => {
                tracing::error!("Error committing & flushing segment writer: {:?}", e);
                terminate_with_error(self.result_channel.take(), Box::new(e), ctx);
            }
        }
    }
}

#[async_trait]
impl Handler<TaskResult<FlushSegmentWriterOutput, FlushSegmentWriterOperatorError>>
    for CompactOrchestrator
{
    type Result = ();

    async fn handle(
        &mut self,
        message: TaskResult<FlushSegmentWriterOutput, FlushSegmentWriterOperatorError>,
        ctx: &crate::system::ComponentContext<CompactOrchestrator>,
    ) {
        let message = message.into_inner();
        match message {
            Ok(message) => {
                // Drops the span so that the end timestamp is accurate
                let _ = self.segment_spans.remove(&message.flush_info.segment_id);

                self.flush_results.push(message.flush_info);
                self.num_uncompleted_flush_tasks -= 1;

                if self.num_uncompleted_flush_tasks == 0 {
                    // Unwrap should be safe here as we are guaranteed to have a value by construction
                    self.register(self.pulled_log_offset.unwrap(), ctx.receiver())
                        .await;
                }
            }
            Err(e) => {
                tracing::error!("Error committing & flushing segment writer: {:?}", e);
                terminate_with_error(self.result_channel.take(), Box::new(e), ctx);
            }
        }
    }
}

#[async_trait]
impl Handler<TaskResult<RegisterOutput, RegisterError>> for CompactOrchestrator {
    type Result = ();

    async fn handle(
        &mut self,
        message: TaskResult<RegisterOutput, RegisterError>,
        ctx: &crate::system::ComponentContext<CompactOrchestrator>,
    ) {
        let message = message.into_inner();
        // Return execution state to the compaction manager
        let result_channel = self
            .result_channel
            .take()
            .expect("Invariant violation. Result channel is not set.");

        match message {
            Ok(_) => {
                let response = CompactionResponse {
                    id: self.id,
                    compaction_job: self.compaction_job.clone(),
                    message: "Compaction Complete".to_string(),
                };
                let _ = result_channel.send(Ok(response));
            }
            Err(e) => {
                tracing::error!("Error registering compaction: {:?}", e);
                terminate_with_error(Some(result_channel), Box::new(e), ctx);
            }
        }
    }
}
