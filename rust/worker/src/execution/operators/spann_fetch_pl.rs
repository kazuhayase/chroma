use chroma_error::{ChromaError, ErrorCodes};
use thiserror::Error;
use tonic::async_trait;

use crate::{
    execution::operator::Operator,
    segment::spann_segment::{SpannSegmentReader, SpannSegmentReaderContext},
};

#[derive(Debug)]
pub struct SpannFetchPlInput {
    reader_context: SpannSegmentReaderContext,
    head_id: u32,
}

#[derive(Debug)]
pub struct SpannFetchPlOutput {
    doc_offset_ids: Vec<u32>,
    doc_versions: Vec<u32>,
    doc_embeddings: Vec<f32>,
}

#[derive(Error, Debug)]
pub enum SpannFetchPlError {
    #[error("Error creating spann segment reader")]
    SpannSegmentReaderCreationError,
    #[error("Error querying reader")]
    SpannSegmentReaderError,
}

impl ChromaError for SpannFetchPlError {
    fn code(&self) -> ErrorCodes {
        match self {
            Self::SpannSegmentReaderCreationError => ErrorCodes::Internal,
            Self::SpannSegmentReaderError => ErrorCodes::Internal,
        }
    }
}

#[derive(Debug)]
pub struct SpannFetchPlOperator {}

impl SpannFetchPlOperator {
    pub fn new() -> Box<Self> {
        Box::new(SpannFetchPlOperator {})
    }
}

#[async_trait]
impl Operator<SpannFetchPlInput, SpannFetchPlOutput> for SpannFetchPlOperator {
    type Error = SpannFetchPlError;

    async fn run(
        &self,
        input: &SpannFetchPlInput,
    ) -> Result<SpannFetchPlOutput, SpannFetchPlError> {
        let spann_reader = SpannSegmentReader::from_segment(
            &input.reader_context.segment,
            &input.reader_context.blockfile_provider,
            &input.reader_context.hnsw_provider,
            input.reader_context.dimension,
        )
        .await
        .map_err(|_| SpannFetchPlError::SpannSegmentReaderCreationError)?;
        let pl = spann_reader
            .fetch_posting_list(input.head_id)
            .await
            .map_err(|_| SpannFetchPlError::SpannSegmentReaderError)?;
        Ok(SpannFetchPlOutput {
            doc_offset_ids: pl.doc_offset_ids.to_vec(),
            doc_versions: pl.doc_versions.to_vec(),
            doc_embeddings: pl.doc_embeddings.to_vec(),
        })
    }
}
