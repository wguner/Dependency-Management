// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: ServerUtils.proto

package ServerUtils

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

// ServerUtilsServicesClient is the client API for ServerUtilsServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServerUtilsServicesClient interface {
	Ping(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (*ServerInfo, error)
}

type serverUtilsServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewServerUtilsServicesClient(cc grpc.ClientConnInterface) ServerUtilsServicesClient {
	return &serverUtilsServicesClient{cc}
}

func (c *serverUtilsServicesClient) Ping(ctx context.Context, in *ClientInfo, opts ...grpc.CallOption) (*ServerInfo, error) {
	out := new(ServerInfo)
	err := c.cc.Invoke(ctx, "/ServerUtils.ServerUtilsServices/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerUtilsServicesServer is the server API for ServerUtilsServices service.
// All implementations must embed UnimplementedServerUtilsServicesServer
// for forward compatibility
type ServerUtilsServicesServer interface {
	Ping(context.Context, *ClientInfo) (*ServerInfo, error)
	mustEmbedUnimplementedServerUtilsServicesServer()
}

// UnimplementedServerUtilsServicesServer must be embedded to have forward compatible implementations.
type UnimplementedServerUtilsServicesServer struct {
}

func (UnimplementedServerUtilsServicesServer) Ping(context.Context, *ClientInfo) (*ServerInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedServerUtilsServicesServer) mustEmbedUnimplementedServerUtilsServicesServer() {}

// UnsafeServerUtilsServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServerUtilsServicesServer will
// result in compilation errors.
type UnsafeServerUtilsServicesServer interface {
	mustEmbedUnimplementedServerUtilsServicesServer()
}

func RegisterServerUtilsServicesServer(s grpc.ServiceRegistrar, srv ServerUtilsServicesServer) {
	s.RegisterService(&ServerUtilsServices_ServiceDesc, srv)
}

func _ServerUtilsServices_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerUtilsServicesServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ServerUtils.ServerUtilsServices/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerUtilsServicesServer).Ping(ctx, req.(*ClientInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// ServerUtilsServices_ServiceDesc is the grpc.ServiceDesc for ServerUtilsServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServerUtilsServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ServerUtils.ServerUtilsServices",
	HandlerType: (*ServerUtilsServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _ServerUtilsServices_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ServerUtils.proto",
}
