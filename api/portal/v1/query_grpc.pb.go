// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: noble/dollar/portal/v1/query.proto

package portalv1

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
	Query_Owner_FullMethodName = "/noble.dollar.portal.v1.Query/Owner"
	Query_Peers_FullMethodName = "/noble.dollar.portal.v1.Query/Peers"
	Query_Nonce_FullMethodName = "/noble.dollar.portal.v1.Query/Nonce"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	Owner(ctx context.Context, in *QueryOwner, opts ...grpc.CallOption) (*QueryOwnerResponse, error)
	Peers(ctx context.Context, in *QueryPeers, opts ...grpc.CallOption) (*QueryPeersResponse, error)
	Nonce(ctx context.Context, in *QueryNonce, opts ...grpc.CallOption) (*QueryNonceResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Owner(ctx context.Context, in *QueryOwner, opts ...grpc.CallOption) (*QueryOwnerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryOwnerResponse)
	err := c.cc.Invoke(ctx, Query_Owner_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Peers(ctx context.Context, in *QueryPeers, opts ...grpc.CallOption) (*QueryPeersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryPeersResponse)
	err := c.cc.Invoke(ctx, Query_Peers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Nonce(ctx context.Context, in *QueryNonce, opts ...grpc.CallOption) (*QueryNonceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryNonceResponse)
	err := c.cc.Invoke(ctx, Query_Nonce_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility.
type QueryServer interface {
	Owner(context.Context, *QueryOwner) (*QueryOwnerResponse, error)
	Peers(context.Context, *QueryPeers) (*QueryPeersResponse, error)
	Nonce(context.Context, *QueryNonce) (*QueryNonceResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueryServer struct{}

func (UnimplementedQueryServer) Owner(context.Context, *QueryOwner) (*QueryOwnerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Owner not implemented")
}
func (UnimplementedQueryServer) Peers(context.Context, *QueryPeers) (*QueryPeersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Peers not implemented")
}
func (UnimplementedQueryServer) Nonce(context.Context, *QueryNonce) (*QueryNonceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Nonce not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}
func (UnimplementedQueryServer) testEmbeddedByValue()               {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	// If the following call pancis, it indicates UnimplementedQueryServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Owner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOwner)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Owner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Owner_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Owner(ctx, req.(*QueryOwner))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Peers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPeers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Peers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Peers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Peers(ctx, req.(*QueryPeers))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Nonce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryNonce)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Nonce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Nonce_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Nonce(ctx, req.(*QueryNonce))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "noble.dollar.portal.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Owner",
			Handler:    _Query_Owner_Handler,
		},
		{
			MethodName: "Peers",
			Handler:    _Query_Peers_Handler,
		},
		{
			MethodName: "Nonce",
			Handler:    _Query_Nonce_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "noble/dollar/portal/v1/query.proto",
}
