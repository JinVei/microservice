// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/v1/app/reply_service.proto

package app

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

// ReplyCommentServiceClient is the client API for ReplyCommentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReplyCommentServiceClient interface {
	ListCommentPage(ctx context.Context, in *ListCommentPageReq, opts ...grpc.CallOption) (*ListCommentPageResp, error)
	PutComment(ctx context.Context, in *PutCommentReq, opts ...grpc.CallOption) (*PutCommentResp, error)
	CreateSubject(ctx context.Context, in *CreateSubjectReq, opts ...grpc.CallOption) (*CreateSubjectResp, error)
	GetSubject(ctx context.Context, in *GetSubjectReq, opts ...grpc.CallOption) (*GetSubjectResp, error)
}

type replyCommentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReplyCommentServiceClient(cc grpc.ClientConnInterface) ReplyCommentServiceClient {
	return &replyCommentServiceClient{cc}
}

func (c *replyCommentServiceClient) ListCommentPage(ctx context.Context, in *ListCommentPageReq, opts ...grpc.CallOption) (*ListCommentPageResp, error) {
	out := new(ListCommentPageResp)
	err := c.cc.Invoke(ctx, "/jv.microservice.v1.app.ReplyCommentService/ListCommentPage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replyCommentServiceClient) PutComment(ctx context.Context, in *PutCommentReq, opts ...grpc.CallOption) (*PutCommentResp, error) {
	out := new(PutCommentResp)
	err := c.cc.Invoke(ctx, "/jv.microservice.v1.app.ReplyCommentService/PutComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replyCommentServiceClient) CreateSubject(ctx context.Context, in *CreateSubjectReq, opts ...grpc.CallOption) (*CreateSubjectResp, error) {
	out := new(CreateSubjectResp)
	err := c.cc.Invoke(ctx, "/jv.microservice.v1.app.ReplyCommentService/CreateSubject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replyCommentServiceClient) GetSubject(ctx context.Context, in *GetSubjectReq, opts ...grpc.CallOption) (*GetSubjectResp, error) {
	out := new(GetSubjectResp)
	err := c.cc.Invoke(ctx, "/jv.microservice.v1.app.ReplyCommentService/GetSubject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplyCommentServiceServer is the server API for ReplyCommentService service.
// All implementations should embed UnimplementedReplyCommentServiceServer
// for forward compatibility
type ReplyCommentServiceServer interface {
	ListCommentPage(context.Context, *ListCommentPageReq) (*ListCommentPageResp, error)
	PutComment(context.Context, *PutCommentReq) (*PutCommentResp, error)
	CreateSubject(context.Context, *CreateSubjectReq) (*CreateSubjectResp, error)
	GetSubject(context.Context, *GetSubjectReq) (*GetSubjectResp, error)
}

// UnimplementedReplyCommentServiceServer should be embedded to have forward compatible implementations.
type UnimplementedReplyCommentServiceServer struct {
}

func (UnimplementedReplyCommentServiceServer) ListCommentPage(context.Context, *ListCommentPageReq) (*ListCommentPageResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCommentPage not implemented")
}
func (UnimplementedReplyCommentServiceServer) PutComment(context.Context, *PutCommentReq) (*PutCommentResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutComment not implemented")
}
func (UnimplementedReplyCommentServiceServer) CreateSubject(context.Context, *CreateSubjectReq) (*CreateSubjectResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSubject not implemented")
}
func (UnimplementedReplyCommentServiceServer) GetSubject(context.Context, *GetSubjectReq) (*GetSubjectResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSubject not implemented")
}

// UnsafeReplyCommentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReplyCommentServiceServer will
// result in compilation errors.
type UnsafeReplyCommentServiceServer interface {
	mustEmbedUnimplementedReplyCommentServiceServer()
}

func RegisterReplyCommentServiceServer(s grpc.ServiceRegistrar, srv ReplyCommentServiceServer) {
	s.RegisterService(&ReplyCommentService_ServiceDesc, srv)
}

func _ReplyCommentService_ListCommentPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCommentPageReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplyCommentServiceServer).ListCommentPage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jv.microservice.v1.app.ReplyCommentService/ListCommentPage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplyCommentServiceServer).ListCommentPage(ctx, req.(*ListCommentPageReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReplyCommentService_PutComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutCommentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplyCommentServiceServer).PutComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jv.microservice.v1.app.ReplyCommentService/PutComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplyCommentServiceServer).PutComment(ctx, req.(*PutCommentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReplyCommentService_CreateSubject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSubjectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplyCommentServiceServer).CreateSubject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jv.microservice.v1.app.ReplyCommentService/CreateSubject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplyCommentServiceServer).CreateSubject(ctx, req.(*CreateSubjectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReplyCommentService_GetSubject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSubjectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplyCommentServiceServer).GetSubject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jv.microservice.v1.app.ReplyCommentService/GetSubject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplyCommentServiceServer).GetSubject(ctx, req.(*GetSubjectReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ReplyCommentService_ServiceDesc is the grpc.ServiceDesc for ReplyCommentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReplyCommentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "jv.microservice.v1.app.ReplyCommentService",
	HandlerType: (*ReplyCommentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListCommentPage",
			Handler:    _ReplyCommentService_ListCommentPage_Handler,
		},
		{
			MethodName: "PutComment",
			Handler:    _ReplyCommentService_PutComment_Handler,
		},
		{
			MethodName: "CreateSubject",
			Handler:    _ReplyCommentService_CreateSubject_Handler,
		},
		{
			MethodName: "GetSubject",
			Handler:    _ReplyCommentService_GetSubject_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/app/reply_service.proto",
}