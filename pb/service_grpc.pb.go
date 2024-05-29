// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: service.proto

package pb

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

// TranscodeClient is the client API for Transcode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TranscodeClient interface {
	BeginOneFile(ctx context.Context, in *BeginOneFileRequest, opts ...grpc.CallOption) (*BeginOneFileReply, error)
	CheckOneFile(ctx context.Context, in *CheckOneFileRequest, opts ...grpc.CallOption) (*CheckOneFileReply, error)
}

type transcodeClient struct {
	cc grpc.ClientConnInterface
}

func NewTranscodeClient(cc grpc.ClientConnInterface) TranscodeClient {
	return &transcodeClient{cc}
}

func (c *transcodeClient) BeginOneFile(ctx context.Context, in *BeginOneFileRequest, opts ...grpc.CallOption) (*BeginOneFileReply, error) {
	out := new(BeginOneFileReply)
	err := c.cc.Invoke(ctx, "/Transcode/BeginOneFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transcodeClient) CheckOneFile(ctx context.Context, in *CheckOneFileRequest, opts ...grpc.CallOption) (*CheckOneFileReply, error) {
	out := new(CheckOneFileReply)
	err := c.cc.Invoke(ctx, "/Transcode/CheckOneFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TranscodeServer is the server API for Transcode service.
// All implementations must embed UnimplementedTranscodeServer
// for forward compatibility
type TranscodeServer interface {
	BeginOneFile(context.Context, *BeginOneFileRequest) (*BeginOneFileReply, error)
	CheckOneFile(context.Context, *CheckOneFileRequest) (*CheckOneFileReply, error)
	mustEmbedUnimplementedTranscodeServer()
}

// UnimplementedTranscodeServer must be embedded to have forward compatible implementations.
type UnimplementedTranscodeServer struct {
}

func (UnimplementedTranscodeServer) BeginOneFile(context.Context, *BeginOneFileRequest) (*BeginOneFileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeginOneFile not implemented")
}
func (UnimplementedTranscodeServer) CheckOneFile(context.Context, *CheckOneFileRequest) (*CheckOneFileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckOneFile not implemented")
}
func (UnimplementedTranscodeServer) mustEmbedUnimplementedTranscodeServer() {}

// UnsafeTranscodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TranscodeServer will
// result in compilation errors.
type UnsafeTranscodeServer interface {
	mustEmbedUnimplementedTranscodeServer()
}

func RegisterTranscodeServer(s grpc.ServiceRegistrar, srv TranscodeServer) {
	s.RegisterService(&Transcode_ServiceDesc, srv)
}

func _Transcode_BeginOneFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BeginOneFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TranscodeServer).BeginOneFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Transcode/BeginOneFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TranscodeServer).BeginOneFile(ctx, req.(*BeginOneFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transcode_CheckOneFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckOneFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TranscodeServer).CheckOneFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Transcode/CheckOneFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TranscodeServer).CheckOneFile(ctx, req.(*CheckOneFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Transcode_ServiceDesc is the grpc.ServiceDesc for Transcode service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Transcode_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Transcode",
	HandlerType: (*TranscodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BeginOneFile",
			Handler:    _Transcode_BeginOneFile_Handler,
		},
		{
			MethodName: "CheckOneFile",
			Handler:    _Transcode_CheckOneFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
