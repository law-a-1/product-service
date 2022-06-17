// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: grpc/order.proto

package grpc

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

// OrderClient is the client API for Order service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderClient interface {
	AddOrder(ctx context.Context, in *AddOrderRequest, opts ...grpc.CallOption) (*AddOrderResponse, error)
}

type orderClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderClient(cc grpc.ClientConnInterface) OrderClient {
	return &orderClient{cc}
}

func (c *orderClient) AddOrder(ctx context.Context, in *AddOrderRequest, opts ...grpc.CallOption) (*AddOrderResponse, error) {
	out := new(AddOrderResponse)
	err := c.cc.Invoke(ctx, "/Order/AddOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServer is the server API for Order service.
// All implementations must embed UnimplementedOrderServer
// for forward compatibility
type OrderServer interface {
	AddOrder(context.Context, *AddOrderRequest) (*AddOrderResponse, error)
	mustEmbedUnimplementedOrderServer()
}

// UnimplementedOrderServer must be embedded to have forward compatible implementations.
type UnimplementedOrderServer struct {
}

func (UnimplementedOrderServer) AddOrder(context.Context, *AddOrderRequest) (*AddOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddOrder not implemented")
}
func (UnimplementedOrderServer) mustEmbedUnimplementedOrderServer() {}

// UnsafeOrderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServer will
// result in compilation errors.
type UnsafeOrderServer interface {
	mustEmbedUnimplementedOrderServer()
}

func RegisterOrderServer(s grpc.ServiceRegistrar, srv OrderServer) {
	s.RegisterService(&Order_ServiceDesc, srv)
}

func _Order_AddOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).AddOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Order/AddOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).AddOrder(ctx, req.(*AddOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Order_ServiceDesc is the grpc.ServiceDesc for Order service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Order_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Order",
	HandlerType: (*OrderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddOrder",
			Handler:    _Order_AddOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/order.proto",
}
