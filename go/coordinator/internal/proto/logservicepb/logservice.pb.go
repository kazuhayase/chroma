// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.20.3
// source: chromadb/proto/logservice.proto

package logservicepb

import (
	coordinatorpb "github.com/chroma/chroma-coordinator/internal/proto/coordinatorpb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PushLogsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectionId string                                 `protobuf:"bytes,1,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty"`
	Records      []*coordinatorpb.SubmitEmbeddingRecord `protobuf:"bytes,2,rep,name=records,proto3" json:"records,omitempty"`
}

func (x *PushLogsRequest) Reset() {
	*x = PushLogsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chromadb_proto_logservice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushLogsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushLogsRequest) ProtoMessage() {}

func (x *PushLogsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chromadb_proto_logservice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushLogsRequest.ProtoReflect.Descriptor instead.
func (*PushLogsRequest) Descriptor() ([]byte, []int) {
	return file_chromadb_proto_logservice_proto_rawDescGZIP(), []int{0}
}

func (x *PushLogsRequest) GetCollectionId() string {
	if x != nil {
		return x.CollectionId
	}
	return ""
}

func (x *PushLogsRequest) GetRecords() []*coordinatorpb.SubmitEmbeddingRecord {
	if x != nil {
		return x.Records
	}
	return nil
}

type PushLogsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RecordCount int32 `protobuf:"varint,1,opt,name=record_count,json=recordCount,proto3" json:"record_count,omitempty"`
}

func (x *PushLogsResponse) Reset() {
	*x = PushLogsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chromadb_proto_logservice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushLogsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushLogsResponse) ProtoMessage() {}

func (x *PushLogsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chromadb_proto_logservice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushLogsResponse.ProtoReflect.Descriptor instead.
func (*PushLogsResponse) Descriptor() ([]byte, []int) {
	return file_chromadb_proto_logservice_proto_rawDescGZIP(), []int{1}
}

func (x *PushLogsResponse) GetRecordCount() int32 {
	if x != nil {
		return x.RecordCount
	}
	return 0
}

type PullLogsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectionId string `protobuf:"bytes,1,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty"`
	StartFromId  int64  `protobuf:"varint,2,opt,name=start_from_id,json=startFromId,proto3" json:"start_from_id,omitempty"`
	BatchSize    int32  `protobuf:"varint,3,opt,name=batch_size,json=batchSize,proto3" json:"batch_size,omitempty"`
}

func (x *PullLogsRequest) Reset() {
	*x = PullLogsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chromadb_proto_logservice_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullLogsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullLogsRequest) ProtoMessage() {}

func (x *PullLogsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chromadb_proto_logservice_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullLogsRequest.ProtoReflect.Descriptor instead.
func (*PullLogsRequest) Descriptor() ([]byte, []int) {
	return file_chromadb_proto_logservice_proto_rawDescGZIP(), []int{2}
}

func (x *PullLogsRequest) GetCollectionId() string {
	if x != nil {
		return x.CollectionId
	}
	return ""
}

func (x *PullLogsRequest) GetStartFromId() int64 {
	if x != nil {
		return x.StartFromId
	}
	return 0
}

func (x *PullLogsRequest) GetBatchSize() int32 {
	if x != nil {
		return x.BatchSize
	}
	return 0
}

type PullLogsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Records []*coordinatorpb.SubmitEmbeddingRecord `protobuf:"bytes,1,rep,name=records,proto3" json:"records,omitempty"`
}

func (x *PullLogsResponse) Reset() {
	*x = PullLogsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chromadb_proto_logservice_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullLogsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullLogsResponse) ProtoMessage() {}

func (x *PullLogsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chromadb_proto_logservice_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullLogsResponse.ProtoReflect.Descriptor instead.
func (*PullLogsResponse) Descriptor() ([]byte, []int) {
	return file_chromadb_proto_logservice_proto_rawDescGZIP(), []int{3}
}

func (x *PullLogsResponse) GetRecords() []*coordinatorpb.SubmitEmbeddingRecord {
	if x != nil {
		return x.Records
	}
	return nil
}

var File_chromadb_proto_logservice_proto protoreflect.FileDescriptor

var file_chromadb_proto_logservice_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x61, 0x64, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x06, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x61, 0x1a, 0x1b, 0x63, 0x68, 0x72, 0x6f, 0x6d,
	0x61, 0x64, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x61,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6f, 0x0a, 0x0f, 0x50, 0x75, 0x73, 0x68, 0x4c, 0x6f,
	0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x37,
	0x0a, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x61, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x45,
	0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x07,
	0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x22, 0x35, 0x0a, 0x10, 0x50, 0x75, 0x73, 0x68, 0x4c,
	0x6f, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x79,
	0x0a, 0x0f, 0x50, 0x75, 0x6c, 0x6c, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x61,
	0x74, 0x63, 0x68, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x62, 0x61, 0x74, 0x63, 0x68, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x4b, 0x0a, 0x10, 0x50, 0x75, 0x6c,
	0x6c, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a,
	0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x61, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x45, 0x6d,
	0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x72,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x32, 0x8e, 0x01, 0x0a, 0x0a, 0x4c, 0x6f, 0x67, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3f, 0x0a, 0x08, 0x50, 0x75, 0x73, 0x68, 0x4c, 0x6f, 0x67,
	0x73, 0x12, 0x17, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x61, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4c,
	0x6f, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x63, 0x68, 0x72,
	0x6f, 0x6d, 0x61, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x08, 0x50, 0x75, 0x6c, 0x6c, 0x4c, 0x6f,
	0x67, 0x73, 0x12, 0x17, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x61, 0x2e, 0x50, 0x75, 0x6c, 0x6c,
	0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x63, 0x68,
	0x72, 0x6f, 0x6d, 0x61, 0x2e, 0x50, 0x75, 0x6c, 0x6c, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x42, 0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x61, 0x2f, 0x63, 0x68, 0x72,
	0x6f, 0x6d, 0x61, 0x2d, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c,
	0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_chromadb_proto_logservice_proto_rawDescOnce sync.Once
	file_chromadb_proto_logservice_proto_rawDescData = file_chromadb_proto_logservice_proto_rawDesc
)

func file_chromadb_proto_logservice_proto_rawDescGZIP() []byte {
	file_chromadb_proto_logservice_proto_rawDescOnce.Do(func() {
		file_chromadb_proto_logservice_proto_rawDescData = protoimpl.X.CompressGZIP(file_chromadb_proto_logservice_proto_rawDescData)
	})
	return file_chromadb_proto_logservice_proto_rawDescData
}

var file_chromadb_proto_logservice_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_chromadb_proto_logservice_proto_goTypes = []interface{}{
	(*PushLogsRequest)(nil),                     // 0: chroma.PushLogsRequest
	(*PushLogsResponse)(nil),                    // 1: chroma.PushLogsResponse
	(*PullLogsRequest)(nil),                     // 2: chroma.PullLogsRequest
	(*PullLogsResponse)(nil),                    // 3: chroma.PullLogsResponse
	(*coordinatorpb.SubmitEmbeddingRecord)(nil), // 4: chroma.SubmitEmbeddingRecord
}
var file_chromadb_proto_logservice_proto_depIdxs = []int32{
	4, // 0: chroma.PushLogsRequest.records:type_name -> chroma.SubmitEmbeddingRecord
	4, // 1: chroma.PullLogsResponse.records:type_name -> chroma.SubmitEmbeddingRecord
	0, // 2: chroma.LogService.PushLogs:input_type -> chroma.PushLogsRequest
	2, // 3: chroma.LogService.PullLogs:input_type -> chroma.PullLogsRequest
	1, // 4: chroma.LogService.PushLogs:output_type -> chroma.PushLogsResponse
	3, // 5: chroma.LogService.PullLogs:output_type -> chroma.PullLogsResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_chromadb_proto_logservice_proto_init() }
func file_chromadb_proto_logservice_proto_init() {
	if File_chromadb_proto_logservice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chromadb_proto_logservice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushLogsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chromadb_proto_logservice_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushLogsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chromadb_proto_logservice_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullLogsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chromadb_proto_logservice_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullLogsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_chromadb_proto_logservice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chromadb_proto_logservice_proto_goTypes,
		DependencyIndexes: file_chromadb_proto_logservice_proto_depIdxs,
		MessageInfos:      file_chromadb_proto_logservice_proto_msgTypes,
	}.Build()
	File_chromadb_proto_logservice_proto = out.File
	file_chromadb_proto_logservice_proto_rawDesc = nil
	file_chromadb_proto_logservice_proto_goTypes = nil
	file_chromadb_proto_logservice_proto_depIdxs = nil
}
