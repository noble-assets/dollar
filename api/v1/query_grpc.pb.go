// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: noble/dollar/v1/query.proto

package dollarv1

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
	Query_Index_FullMethodName     = "/noble.dollar.v1.Query/Index"
	Query_Paused_FullMethodName    = "/noble.dollar.v1.Query/Paused"
	Query_Principal_FullMethodName = "/noble.dollar.v1.Query/Principal"
	Query_Yield_FullMethodName     = "/noble.dollar.v1.Query/Yield"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	Index(ctx context.Context, in *QueryIndex, opts ...grpc.CallOption) (*QueryIndexResponse, error)
	Paused(ctx context.Context, in *QueryPaused, opts ...grpc.CallOption) (*QueryPausedResponse, error)
	Principal(ctx context.Context, in *QueryPrincipal, opts ...grpc.CallOption) (*QueryPrincipalResponse, error)
	Yield(ctx context.Context, in *QueryYield, opts ...grpc.CallOption) (*QueryYieldResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Index(ctx context.Context, in *QueryIndex, opts ...grpc.CallOption) (*QueryIndexResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryIndexResponse)
	err := c.cc.Invoke(ctx, Query_Index_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Paused(ctx context.Context, in *QueryPaused, opts ...grpc.CallOption) (*QueryPausedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryPausedResponse)
	err := c.cc.Invoke(ctx, Query_Paused_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Principal(ctx context.Context, in *QueryPrincipal, opts ...grpc.CallOption) (*QueryPrincipalResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryPrincipalResponse)
	err := c.cc.Invoke(ctx, Query_Principal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Yield(ctx context.Context, in *QueryYield, opts ...grpc.CallOption) (*QueryYieldResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryYieldResponse)
	err := c.cc.Invoke(ctx, Query_Yield_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility.
type QueryServer interface {
	Index(context.Context, *QueryIndex) (*QueryIndexResponse, error)
	Paused(context.Context, *QueryPaused) (*QueryPausedResponse, error)
	Principal(context.Context, *QueryPrincipal) (*QueryPrincipalResponse, error)
	Yield(context.Context, *QueryYield) (*QueryYieldResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueryServer struct{}

func (UnimplementedQueryServer) Index(context.Context, *QueryIndex) (*QueryIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}
func (UnimplementedQueryServer) Paused(context.Context, *QueryPaused) (*QueryPausedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Paused not implemented")
}
func (UnimplementedQueryServer) Principal(context.Context, *QueryPrincipal) (*QueryPrincipalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Principal not implemented")
}
func (UnimplementedQueryServer) Yield(context.Context, *QueryYield) (*QueryYieldResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Yield not implemented")
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

func _Query_Index_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIndex)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Index(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Index_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Index(ctx, req.(*QueryIndex))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Paused_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPaused)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Paused(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Paused_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Paused(ctx, req.(*QueryPaused))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Principal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPrincipal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Principal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Principal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Principal(ctx, req.(*QueryPrincipal))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Yield_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryYield)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Yield(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Yield_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Yield(ctx, req.(*QueryYield))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "noble.dollar.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Index",
			Handler:    _Query_Index_Handler,
		},
		{
			MethodName: "Paused",
			Handler:    _Query_Paused_Handler,
		},
		{
			MethodName: "Principal",
			Handler:    _Query_Principal_Handler,
		},
		{
			MethodName: "Yield",
			Handler:    _Query_Yield_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "noble/dollar/v1/query.proto",
}
