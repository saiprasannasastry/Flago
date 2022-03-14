// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: flago.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FlagoServiceClient is the client API for FlagoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FlagoServiceClient interface {
	//Create Flag allows you to define new set of flags
	CreateFlag(ctx context.Context, in *CreateFlagReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	//GetFlag allows the control plane to query dataPlane to
	// validate if the flag is enabled or not
	GetFlag(ctx context.Context, in *FlagReq, opts ...grpc.CallOption) (*FlagResp, error)
	//GetFlags returns list of all flags enabled for customer
	GetFlags(ctx context.Context, in *FlagReq, opts ...grpc.CallOption) (*GetFlagResp, error)
	// OnFlag turns the flag on
	// so the control plane can check get get all data
	// from Data plane when Flag is enabled
	OnFlag(ctx context.Context, in *FlagReq, opts ...grpc.CallOption) (*FlagResp, error)
	// OffFlag turns the flag on
	// so the control plane can check get get all data
	// from Data plane when Flag is enabled
	OffFlag(ctx context.Context, in *FlagReq, opts ...grpc.CallOption) (*FlagResp, error)
}

type flagoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFlagoServiceClient(cc grpc.ClientConnInterface) FlagoServiceClient {
	return &flagoServiceClient{cc}
}

func (c *flagoServiceClient) CreateFlag(ctx context.Context, in *CreateFlagReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/flago.FlagoService/CreateFlag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *flagoServiceClient) GetFlag(ctx context.Context, in *FlagReq, opts ...grpc.CallOption) (*FlagResp, error) {
	out := new(FlagResp)
	err := c.cc.Invoke(ctx, "/flago.FlagoService/GetFlag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *flagoServiceClient) GetFlags(ctx context.Context, in *FlagReq, opts ...grpc.CallOption) (*GetFlagResp, error) {
	out := new(GetFlagResp)
	err := c.cc.Invoke(ctx, "/flago.FlagoService/GetFlags", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *flagoServiceClient) OnFlag(ctx context.Context, in *FlagReq, opts ...grpc.CallOption) (*FlagResp, error) {
	out := new(FlagResp)
	err := c.cc.Invoke(ctx, "/flago.FlagoService/OnFlag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *flagoServiceClient) OffFlag(ctx context.Context, in *FlagReq, opts ...grpc.CallOption) (*FlagResp, error) {
	out := new(FlagResp)
	err := c.cc.Invoke(ctx, "/flago.FlagoService/OffFlag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FlagoServiceServer is the server API for FlagoService service.
// All implementations must embed UnimplementedFlagoServiceServer
// for forward compatibility
type FlagoServiceServer interface {
	//Create Flag allows you to define new set of flags
	CreateFlag(context.Context, *CreateFlagReq) (*emptypb.Empty, error)
	//GetFlag allows the control plane to query dataPlane to
	// validate if the flag is enabled or not
	GetFlag(context.Context, *FlagReq) (*FlagResp, error)
	//GetFlags returns list of all flags enabled for customer
	GetFlags(context.Context, *FlagReq) (*GetFlagResp, error)
	// OnFlag turns the flag on
	// so the control plane can check get get all data
	// from Data plane when Flag is enabled
	OnFlag(context.Context, *FlagReq) (*FlagResp, error)
	// OffFlag turns the flag on
	// so the control plane can check get get all data
	// from Data plane when Flag is enabled
	OffFlag(context.Context, *FlagReq) (*FlagResp, error)
	mustEmbedUnimplementedFlagoServiceServer()
}

// UnimplementedFlagoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFlagoServiceServer struct {
}

func (UnimplementedFlagoServiceServer) CreateFlag(context.Context, *CreateFlagReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFlag not implemented")
}
func (UnimplementedFlagoServiceServer) GetFlag(context.Context, *FlagReq) (*FlagResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFlag not implemented")
}
func (UnimplementedFlagoServiceServer) GetFlags(context.Context, *FlagReq) (*GetFlagResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFlags not implemented")
}
func (UnimplementedFlagoServiceServer) OnFlag(context.Context, *FlagReq) (*FlagResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnFlag not implemented")
}
func (UnimplementedFlagoServiceServer) OffFlag(context.Context, *FlagReq) (*FlagResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OffFlag not implemented")
}
func (UnimplementedFlagoServiceServer) mustEmbedUnimplementedFlagoServiceServer() {}

// UnsafeFlagoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FlagoServiceServer will
// result in compilation errors.
type UnsafeFlagoServiceServer interface {
	mustEmbedUnimplementedFlagoServiceServer()
}

func RegisterFlagoServiceServer(s grpc.ServiceRegistrar, srv FlagoServiceServer) {
	s.RegisterService(&FlagoService_ServiceDesc, srv)
}

func _FlagoService_CreateFlag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFlagReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlagoServiceServer).CreateFlag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flago.FlagoService/CreateFlag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlagoServiceServer).CreateFlag(ctx, req.(*CreateFlagReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FlagoService_GetFlag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FlagReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlagoServiceServer).GetFlag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flago.FlagoService/GetFlag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlagoServiceServer).GetFlag(ctx, req.(*FlagReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FlagoService_GetFlags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FlagReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlagoServiceServer).GetFlags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flago.FlagoService/GetFlags",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlagoServiceServer).GetFlags(ctx, req.(*FlagReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FlagoService_OnFlag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FlagReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlagoServiceServer).OnFlag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flago.FlagoService/OnFlag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlagoServiceServer).OnFlag(ctx, req.(*FlagReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FlagoService_OffFlag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FlagReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlagoServiceServer).OffFlag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flago.FlagoService/OffFlag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlagoServiceServer).OffFlag(ctx, req.(*FlagReq))
	}
	return interceptor(ctx, in, info, handler)
}

// FlagoService_ServiceDesc is the grpc.ServiceDesc for FlagoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FlagoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "flago.FlagoService",
	HandlerType: (*FlagoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFlag",
			Handler:    _FlagoService_CreateFlag_Handler,
		},
		{
			MethodName: "GetFlag",
			Handler:    _FlagoService_GetFlag_Handler,
		},
		{
			MethodName: "GetFlags",
			Handler:    _FlagoService_GetFlags_Handler,
		},
		{
			MethodName: "OnFlag",
			Handler:    _FlagoService_OnFlag_Handler,
		},
		{
			MethodName: "OffFlag",
			Handler:    _FlagoService_OffFlag_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "flago.proto",
}
