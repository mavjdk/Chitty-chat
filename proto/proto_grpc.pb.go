// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: proto.proto

package proto

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

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	MessageRoute(ctx context.Context, opts ...grpc.CallOption) (MessageService_MessageRouteClient, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) MessageRoute(ctx context.Context, opts ...grpc.CallOption) (MessageService_MessageRouteClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessageService_ServiceDesc.Streams[0], "/chittychat.messageService/MessageRoute", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageServiceMessageRouteClient{stream}
	return x, nil
}

type MessageService_MessageRouteClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type messageServiceMessageRouteClient struct {
	grpc.ClientStream
}

func (x *messageServiceMessageRouteClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messageServiceMessageRouteClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations must embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	MessageRoute(MessageService_MessageRouteServer) error
	mustEmbedUnimplementedMessageServiceServer()
}

// UnimplementedMessageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) MessageRoute(MessageService_MessageRouteServer) error {
	return status.Errorf(codes.Unimplemented, "method MessageRoute not implemented")
}
func (UnimplementedMessageServiceServer) mustEmbedUnimplementedMessageServiceServer() {}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_MessageRoute_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessageServiceServer).MessageRoute(&messageServiceMessageRouteServer{stream})
}

type MessageService_MessageRouteServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type messageServiceMessageRouteServer struct {
	grpc.ServerStream
}

func (x *messageServiceMessageRouteServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messageServiceMessageRouteServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chittychat.messageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MessageRoute",
			Handler:       _MessageService_MessageRoute_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto.proto",
}
