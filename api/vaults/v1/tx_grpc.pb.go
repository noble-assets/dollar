// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: noble/dollar/vaults/v1/tx.proto

package vaultsv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Msg_Lock_FullMethodName           = "/noble.dollar.vaults.v1.Msg/Lock"
	Msg_Unlock_FullMethodName         = "/noble.dollar.vaults.v1.Msg/Unlock"
	Msg_SetPausedState_FullMethodName = "/noble.dollar.vaults.v1.Msg/SetPausedState"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	Lock(ctx context.Context, in *MsgLock, opts ...grpc.CallOption) (*MsgLockResponse, error)
	Unlock(ctx context.Context, in *MsgUnlock, opts ...grpc.CallOption) (*MsgUnlockResponse, error)
	SetPausedState(ctx context.Context, in *MsgSetPausedState, opts ...grpc.CallOption) (*MsgSetPausedStateResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) Lock(ctx context.Context, in *MsgLock, opts ...grpc.CallOption) (*MsgLockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgLockResponse)
	err := c.cc.Invoke(ctx, Msg_Lock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Unlock(ctx context.Context, in *MsgUnlock, opts ...grpc.CallOption) (*MsgUnlockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgUnlockResponse)
	err := c.cc.Invoke(ctx, Msg_Unlock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetPausedState(ctx context.Context, in *MsgSetPausedState, opts ...grpc.CallOption) (*MsgSetPausedStateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgSetPausedStateResponse)
	err := c.cc.Invoke(ctx, Msg_SetPausedState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility.
type MsgServer interface {
	Lock(context.Context, *MsgLock) (*MsgLockResponse, error)
	Unlock(context.Context, *MsgUnlock) (*MsgUnlockResponse, error)
	SetPausedState(context.Context, *MsgSetPausedState) (*MsgSetPausedStateResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMsgServer struct{}

func (UnimplementedMsgServer) Lock(context.Context, *MsgLock) (*MsgLockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lock not implemented")
}
func (UnimplementedMsgServer) Unlock(context.Context, *MsgUnlock) (*MsgUnlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unlock not implemented")
}
func (UnimplementedMsgServer) SetPausedState(context.Context, *MsgSetPausedState) (*MsgSetPausedStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPausedState not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}
func (UnimplementedMsgServer) testEmbeddedByValue()             {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	// If the following call pancis, it indicates UnimplementedMsgServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_Lock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgLock)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Lock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Lock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Lock(ctx, req.(*MsgLock))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Unlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUnlock)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Unlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Unlock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Unlock(ctx, req.(*MsgUnlock))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetPausedState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetPausedState)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetPausedState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_SetPausedState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetPausedState(ctx, req.(*MsgSetPausedState))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "noble.dollar.vaults.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Lock",
			Handler:    _Msg_Lock_Handler,
		},
		{
			MethodName: "Unlock",
			Handler:    _Msg_Unlock_Handler,
		},
		{
			MethodName: "SetPausedState",
			Handler:    _Msg_SetPausedState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "noble/dollar/vaults/v1/tx.proto",
}
