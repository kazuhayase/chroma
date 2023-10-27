// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: chromadb/proto/chroma.proto

package coordinatorpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	SegmentServer_LoadSegment_FullMethodName    = "/chroma.SegmentServer/LoadSegment"
	SegmentServer_ReleaseSegment_FullMethodName = "/chroma.SegmentServer/ReleaseSegment"
)

// SegmentServerClient is the client API for SegmentServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SegmentServerClient interface {
	LoadSegment(ctx context.Context, in *LoadSegmentRequest, opts ...grpc.CallOption) (*LoadSegmentResponse, error)
	ReleaseSegment(ctx context.Context, in *ReleaseSegmentRequest, opts ...grpc.CallOption) (*ReleaseSegmentResponse, error)
}

type segmentServerClient struct {
	cc grpc.ClientConnInterface
}

func NewSegmentServerClient(cc grpc.ClientConnInterface) SegmentServerClient {
	return &segmentServerClient{cc}
}

func (c *segmentServerClient) LoadSegment(ctx context.Context, in *LoadSegmentRequest, opts ...grpc.CallOption) (*LoadSegmentResponse, error) {
	out := new(LoadSegmentResponse)
	err := c.cc.Invoke(ctx, SegmentServer_LoadSegment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *segmentServerClient) ReleaseSegment(ctx context.Context, in *ReleaseSegmentRequest, opts ...grpc.CallOption) (*ReleaseSegmentResponse, error) {
	out := new(ReleaseSegmentResponse)
	err := c.cc.Invoke(ctx, SegmentServer_ReleaseSegment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SegmentServerServer is the server API for SegmentServer service.
// All implementations must embed UnimplementedSegmentServerServer
// for forward compatibility
type SegmentServerServer interface {
	LoadSegment(context.Context, *LoadSegmentRequest) (*LoadSegmentResponse, error)
	ReleaseSegment(context.Context, *ReleaseSegmentRequest) (*ReleaseSegmentResponse, error)
	mustEmbedUnimplementedSegmentServerServer()
}

// UnimplementedSegmentServerServer must be embedded to have forward compatible implementations.
type UnimplementedSegmentServerServer struct {
}

func (UnimplementedSegmentServerServer) LoadSegment(context.Context, *LoadSegmentRequest) (*LoadSegmentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadSegment not implemented")
}
func (UnimplementedSegmentServerServer) ReleaseSegment(context.Context, *ReleaseSegmentRequest) (*ReleaseSegmentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseSegment not implemented")
}
func (UnimplementedSegmentServerServer) mustEmbedUnimplementedSegmentServerServer() {}

// UnsafeSegmentServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SegmentServerServer will
// result in compilation errors.
type UnsafeSegmentServerServer interface {
	mustEmbedUnimplementedSegmentServerServer()
}

func RegisterSegmentServerServer(s grpc.ServiceRegistrar, srv SegmentServerServer) {
	s.RegisterService(&SegmentServer_ServiceDesc, srv)
}

func _SegmentServer_LoadSegment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadSegmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SegmentServerServer).LoadSegment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SegmentServer_LoadSegment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SegmentServerServer).LoadSegment(ctx, req.(*LoadSegmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SegmentServer_ReleaseSegment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseSegmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SegmentServerServer).ReleaseSegment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SegmentServer_ReleaseSegment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SegmentServerServer).ReleaseSegment(ctx, req.(*ReleaseSegmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SegmentServer_ServiceDesc is the grpc.ServiceDesc for SegmentServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SegmentServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chroma.SegmentServer",
	HandlerType: (*SegmentServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoadSegment",
			Handler:    _SegmentServer_LoadSegment_Handler,
		},
		{
			MethodName: "ReleaseSegment",
			Handler:    _SegmentServer_ReleaseSegment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chromadb/proto/chroma.proto",
}

const (
	VectorReader_GetVectors_FullMethodName   = "/chroma.VectorReader/GetVectors"
	VectorReader_QueryVectors_FullMethodName = "/chroma.VectorReader/QueryVectors"
)

// VectorReaderClient is the client API for VectorReader service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VectorReaderClient interface {
	GetVectors(ctx context.Context, in *GetVectorsRequest, opts ...grpc.CallOption) (*GetVectorsResponse, error)
	QueryVectors(ctx context.Context, in *QueryVectorsRequest, opts ...grpc.CallOption) (*QueryVectorsResponse, error)
}

type vectorReaderClient struct {
	cc grpc.ClientConnInterface
}

func NewVectorReaderClient(cc grpc.ClientConnInterface) VectorReaderClient {
	return &vectorReaderClient{cc}
}

func (c *vectorReaderClient) GetVectors(ctx context.Context, in *GetVectorsRequest, opts ...grpc.CallOption) (*GetVectorsResponse, error) {
	out := new(GetVectorsResponse)
	err := c.cc.Invoke(ctx, VectorReader_GetVectors_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vectorReaderClient) QueryVectors(ctx context.Context, in *QueryVectorsRequest, opts ...grpc.CallOption) (*QueryVectorsResponse, error) {
	out := new(QueryVectorsResponse)
	err := c.cc.Invoke(ctx, VectorReader_QueryVectors_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VectorReaderServer is the server API for VectorReader service.
// All implementations must embed UnimplementedVectorReaderServer
// for forward compatibility
type VectorReaderServer interface {
	GetVectors(context.Context, *GetVectorsRequest) (*GetVectorsResponse, error)
	QueryVectors(context.Context, *QueryVectorsRequest) (*QueryVectorsResponse, error)
	mustEmbedUnimplementedVectorReaderServer()
}

// UnimplementedVectorReaderServer must be embedded to have forward compatible implementations.
type UnimplementedVectorReaderServer struct {
}

func (UnimplementedVectorReaderServer) GetVectors(context.Context, *GetVectorsRequest) (*GetVectorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVectors not implemented")
}
func (UnimplementedVectorReaderServer) QueryVectors(context.Context, *QueryVectorsRequest) (*QueryVectorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryVectors not implemented")
}
func (UnimplementedVectorReaderServer) mustEmbedUnimplementedVectorReaderServer() {}

// UnsafeVectorReaderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VectorReaderServer will
// result in compilation errors.
type UnsafeVectorReaderServer interface {
	mustEmbedUnimplementedVectorReaderServer()
}

func RegisterVectorReaderServer(s grpc.ServiceRegistrar, srv VectorReaderServer) {
	s.RegisterService(&VectorReader_ServiceDesc, srv)
}

func _VectorReader_GetVectors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVectorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VectorReaderServer).GetVectors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VectorReader_GetVectors_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VectorReaderServer).GetVectors(ctx, req.(*GetVectorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VectorReader_QueryVectors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryVectorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VectorReaderServer).QueryVectors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VectorReader_QueryVectors_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VectorReaderServer).QueryVectors(ctx, req.(*QueryVectorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VectorReader_ServiceDesc is the grpc.ServiceDesc for VectorReader service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VectorReader_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chroma.VectorReader",
	HandlerType: (*VectorReaderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVectors",
			Handler:    _VectorReader_GetVectors_Handler,
		},
		{
			MethodName: "QueryVectors",
			Handler:    _VectorReader_QueryVectors_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chromadb/proto/chroma.proto",
}
