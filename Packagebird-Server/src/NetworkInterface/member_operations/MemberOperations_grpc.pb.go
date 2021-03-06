// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: MemberOperations.proto

package member_operations

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

// MemberCRUDServicesClient is the client API for MemberCRUDServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MemberCRUDServicesClient interface {
	CreateMember(ctx context.Context, in *MemberRequest, opts ...grpc.CallOption) (*MemberResponse, error)
	RemoveMember(ctx context.Context, in *MemberRequest, opts ...grpc.CallOption) (*MemberResponse, error)
}

type memberCRUDServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewMemberCRUDServicesClient(cc grpc.ClientConnInterface) MemberCRUDServicesClient {
	return &memberCRUDServicesClient{cc}
}

func (c *memberCRUDServicesClient) CreateMember(ctx context.Context, in *MemberRequest, opts ...grpc.CallOption) (*MemberResponse, error) {
	out := new(MemberResponse)
	err := c.cc.Invoke(ctx, "/MemberOperations.MemberCRUDServices/CreateMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberCRUDServicesClient) RemoveMember(ctx context.Context, in *MemberRequest, opts ...grpc.CallOption) (*MemberResponse, error) {
	out := new(MemberResponse)
	err := c.cc.Invoke(ctx, "/MemberOperations.MemberCRUDServices/RemoveMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MemberCRUDServicesServer is the server API for MemberCRUDServices service.
// All implementations must embed UnimplementedMemberCRUDServicesServer
// for forward compatibility
type MemberCRUDServicesServer interface {
	CreateMember(context.Context, *MemberRequest) (*MemberResponse, error)
	RemoveMember(context.Context, *MemberRequest) (*MemberResponse, error)
	mustEmbedUnimplementedMemberCRUDServicesServer()
}

// UnimplementedMemberCRUDServicesServer must be embedded to have forward compatible implementations.
type UnimplementedMemberCRUDServicesServer struct {
}

func (UnimplementedMemberCRUDServicesServer) CreateMember(context.Context, *MemberRequest) (*MemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMember not implemented")
}
func (UnimplementedMemberCRUDServicesServer) RemoveMember(context.Context, *MemberRequest) (*MemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMember not implemented")
}
func (UnimplementedMemberCRUDServicesServer) mustEmbedUnimplementedMemberCRUDServicesServer() {}

// UnsafeMemberCRUDServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MemberCRUDServicesServer will
// result in compilation errors.
type UnsafeMemberCRUDServicesServer interface {
	mustEmbedUnimplementedMemberCRUDServicesServer()
}

func RegisterMemberCRUDServicesServer(s grpc.ServiceRegistrar, srv MemberCRUDServicesServer) {
	s.RegisterService(&MemberCRUDServices_ServiceDesc, srv)
}

func _MemberCRUDServices_CreateMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberCRUDServicesServer).CreateMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MemberOperations.MemberCRUDServices/CreateMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberCRUDServicesServer).CreateMember(ctx, req.(*MemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MemberCRUDServices_RemoveMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberCRUDServicesServer).RemoveMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MemberOperations.MemberCRUDServices/RemoveMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberCRUDServicesServer).RemoveMember(ctx, req.(*MemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MemberCRUDServices_ServiceDesc is the grpc.ServiceDesc for MemberCRUDServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MemberCRUDServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MemberOperations.MemberCRUDServices",
	HandlerType: (*MemberCRUDServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMember",
			Handler:    _MemberCRUDServices_CreateMember_Handler,
		},
		{
			MethodName: "RemoveMember",
			Handler:    _MemberCRUDServices_RemoveMember_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "MemberOperations.proto",
}

// MemberAuthenticationClient is the client API for MemberAuthentication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MemberAuthenticationClient interface {
	AuthenticateMember(ctx context.Context, in *MemberAuthenticationRequest, opts ...grpc.CallOption) (*MemberAuthenticationResponse, error)
}

type memberAuthenticationClient struct {
	cc grpc.ClientConnInterface
}

func NewMemberAuthenticationClient(cc grpc.ClientConnInterface) MemberAuthenticationClient {
	return &memberAuthenticationClient{cc}
}

func (c *memberAuthenticationClient) AuthenticateMember(ctx context.Context, in *MemberAuthenticationRequest, opts ...grpc.CallOption) (*MemberAuthenticationResponse, error) {
	out := new(MemberAuthenticationResponse)
	err := c.cc.Invoke(ctx, "/MemberOperations.MemberAuthentication/AuthenticateMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MemberAuthenticationServer is the server API for MemberAuthentication service.
// All implementations must embed UnimplementedMemberAuthenticationServer
// for forward compatibility
type MemberAuthenticationServer interface {
	AuthenticateMember(context.Context, *MemberAuthenticationRequest) (*MemberAuthenticationResponse, error)
	mustEmbedUnimplementedMemberAuthenticationServer()
}

// UnimplementedMemberAuthenticationServer must be embedded to have forward compatible implementations.
type UnimplementedMemberAuthenticationServer struct {
}

func (UnimplementedMemberAuthenticationServer) AuthenticateMember(context.Context, *MemberAuthenticationRequest) (*MemberAuthenticationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateMember not implemented")
}
func (UnimplementedMemberAuthenticationServer) mustEmbedUnimplementedMemberAuthenticationServer() {}

// UnsafeMemberAuthenticationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MemberAuthenticationServer will
// result in compilation errors.
type UnsafeMemberAuthenticationServer interface {
	mustEmbedUnimplementedMemberAuthenticationServer()
}

func RegisterMemberAuthenticationServer(s grpc.ServiceRegistrar, srv MemberAuthenticationServer) {
	s.RegisterService(&MemberAuthentication_ServiceDesc, srv)
}

func _MemberAuthentication_AuthenticateMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MemberAuthenticationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberAuthenticationServer).AuthenticateMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MemberOperations.MemberAuthentication/AuthenticateMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberAuthenticationServer).AuthenticateMember(ctx, req.(*MemberAuthenticationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MemberAuthentication_ServiceDesc is the grpc.ServiceDesc for MemberAuthentication service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MemberAuthentication_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MemberOperations.MemberAuthentication",
	HandlerType: (*MemberAuthenticationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthenticateMember",
			Handler:    _MemberAuthentication_AuthenticateMember_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "MemberOperations.proto",
}
