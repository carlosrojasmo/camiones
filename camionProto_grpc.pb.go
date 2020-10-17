// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// OrdenServiceClient is the client API for OrdenService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrdenServiceClient interface {
	GetPack(ctx context.Context, in *AskForPack, opts ...grpc.CallOption) (*SendPack, error)
	Report(ctx context.Context, in *ReportDelivery, opts ...grpc.CallOption) (*ReportOk, error)
}

type ordenServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrdenServiceClient(cc grpc.ClientConnInterface) OrdenServiceClient {
	return &ordenServiceClient{cc}
}

func (c *ordenServiceClient) GetPack(ctx context.Context, in *AskForPack, opts ...grpc.CallOption) (*SendPack, error) {
	out := new(SendPack)
	err := c.cc.Invoke(ctx, "/proto.OrdenService/getPack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ordenServiceClient) Report(ctx context.Context, in *ReportDelivery, opts ...grpc.CallOption) (*ReportOk, error) {
	out := new(ReportOk)
	err := c.cc.Invoke(ctx, "/proto.OrdenService/report", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrdenServiceServer is the server API for OrdenService service.
// All implementations must embed UnimplementedOrdenServiceServer
// for forward compatibility
type OrdenServiceServer interface {
	GetPack(context.Context, *AskForPack) (*SendPack, error)
	Report(context.Context, *ReportDelivery) (*ReportOk, error)
	mustEmbedUnimplementedOrdenServiceServer()
}

// UnimplementedOrdenServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOrdenServiceServer struct {
}

func (UnimplementedOrdenServiceServer) GetPack(context.Context, *AskForPack) (*SendPack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPack not implemented")
}
func (UnimplementedOrdenServiceServer) Report(context.Context, *ReportDelivery) (*ReportOk, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Report not implemented")
}
func (UnimplementedOrdenServiceServer) mustEmbedUnimplementedOrdenServiceServer() {}

// UnsafeOrdenServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrdenServiceServer will
// result in compilation errors.
type UnsafeOrdenServiceServer interface {
	mustEmbedUnimplementedOrdenServiceServer()
}

func RegisterOrdenServiceServer(s *grpc.Server, srv OrdenServiceServer) {
	s.RegisterService(&_OrdenService_serviceDesc, srv)
}

func _OrdenService_GetPack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AskForPack)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdenServiceServer).GetPack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrdenService/GetPack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdenServiceServer).GetPack(ctx, req.(*AskForPack))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrdenService_Report_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportDelivery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdenServiceServer).Report(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrdenService/Report",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdenServiceServer).Report(ctx, req.(*ReportDelivery))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrdenService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.OrdenService",
	HandlerType: (*OrdenServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getPack",
			Handler:    _OrdenService_GetPack_Handler,
		},
		{
			MethodName: "report",
			Handler:    _OrdenService_Report_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "camionProto.proto",
}
