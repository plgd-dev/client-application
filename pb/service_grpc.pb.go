// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// DeviceGatewayClient is the client API for DeviceGateway service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceGatewayClient interface {
	// Discover devices by client application. This operation fills cache of mappings deviceId to endpoints and this cache is used by other RPC calls.
	GetDevices(ctx context.Context, in *GetDevicesRequest, opts ...grpc.CallOption) (DeviceGateway_GetDevicesClient, error)
	// Get device via endpoint address. This operation fills cache of mappings deviceId to endpoints and this cache is used by other RPC calls.
	GetDeviceByEndpoint(ctx context.Context, in *GetDeviceByEndpointRequest, opts ...grpc.CallOption) (*Device, error)
	// Get resource from the device.
	GetResource(ctx context.Context, in *GetResourceRequest, opts ...grpc.CallOption) (*GetResourceResponse, error)
	// Update resource at the device.
	UpdateResource(ctx context.Context, in *UpdateResourceRequest, opts ...grpc.CallOption) (*UpdateResourceResponse, error)
}

type deviceGatewayClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceGatewayClient(cc grpc.ClientConnInterface) DeviceGatewayClient {
	return &deviceGatewayClient{cc}
}

func (c *deviceGatewayClient) GetDevices(ctx context.Context, in *GetDevicesRequest, opts ...grpc.CallOption) (DeviceGateway_GetDevicesClient, error) {
	stream, err := c.cc.NewStream(ctx, &DeviceGateway_ServiceDesc.Streams[0], "/service.pb.DeviceGateway/GetDevices", opts...)
	if err != nil {
		return nil, err
	}
	x := &deviceGatewayGetDevicesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DeviceGateway_GetDevicesClient interface {
	Recv() (*Device, error)
	grpc.ClientStream
}

type deviceGatewayGetDevicesClient struct {
	grpc.ClientStream
}

func (x *deviceGatewayGetDevicesClient) Recv() (*Device, error) {
	m := new(Device)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *deviceGatewayClient) GetDeviceByEndpoint(ctx context.Context, in *GetDeviceByEndpointRequest, opts ...grpc.CallOption) (*Device, error) {
	out := new(Device)
	err := c.cc.Invoke(ctx, "/service.pb.DeviceGateway/GetDeviceByEndpoint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceGatewayClient) GetResource(ctx context.Context, in *GetResourceRequest, opts ...grpc.CallOption) (*GetResourceResponse, error) {
	out := new(GetResourceResponse)
	err := c.cc.Invoke(ctx, "/service.pb.DeviceGateway/GetResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceGatewayClient) UpdateResource(ctx context.Context, in *UpdateResourceRequest, opts ...grpc.CallOption) (*UpdateResourceResponse, error) {
	out := new(UpdateResourceResponse)
	err := c.cc.Invoke(ctx, "/service.pb.DeviceGateway/UpdateResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceGatewayServer is the server API for DeviceGateway service.
// All implementations must embed UnimplementedDeviceGatewayServer
// for forward compatibility
type DeviceGatewayServer interface {
	// Discover devices by client application. This operation fills cache of mappings deviceId to endpoints and this cache is used by other RPC calls.
	GetDevices(*GetDevicesRequest, DeviceGateway_GetDevicesServer) error
	// Get device via endpoint address. This operation fills cache of mappings deviceId to endpoints and this cache is used by other RPC calls.
	GetDeviceByEndpoint(context.Context, *GetDeviceByEndpointRequest) (*Device, error)
	// Get resource from the device.
	GetResource(context.Context, *GetResourceRequest) (*GetResourceResponse, error)
	// Update resource at the device.
	UpdateResource(context.Context, *UpdateResourceRequest) (*UpdateResourceResponse, error)
	mustEmbedUnimplementedDeviceGatewayServer()
}

// UnimplementedDeviceGatewayServer must be embedded to have forward compatible implementations.
type UnimplementedDeviceGatewayServer struct {
}

func (UnimplementedDeviceGatewayServer) GetDevices(*GetDevicesRequest, DeviceGateway_GetDevicesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetDevices not implemented")
}
func (UnimplementedDeviceGatewayServer) GetDeviceByEndpoint(context.Context, *GetDeviceByEndpointRequest) (*Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeviceByEndpoint not implemented")
}
func (UnimplementedDeviceGatewayServer) GetResource(context.Context, *GetResourceRequest) (*GetResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResource not implemented")
}
func (UnimplementedDeviceGatewayServer) UpdateResource(context.Context, *UpdateResourceRequest) (*UpdateResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateResource not implemented")
}
func (UnimplementedDeviceGatewayServer) mustEmbedUnimplementedDeviceGatewayServer() {}

// UnsafeDeviceGatewayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceGatewayServer will
// result in compilation errors.
type UnsafeDeviceGatewayServer interface {
	mustEmbedUnimplementedDeviceGatewayServer()
}

func RegisterDeviceGatewayServer(s grpc.ServiceRegistrar, srv DeviceGatewayServer) {
	s.RegisterService(&DeviceGateway_ServiceDesc, srv)
}

func _DeviceGateway_GetDevices_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetDevicesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DeviceGatewayServer).GetDevices(m, &deviceGatewayGetDevicesServer{stream})
}

type DeviceGateway_GetDevicesServer interface {
	Send(*Device) error
	grpc.ServerStream
}

type deviceGatewayGetDevicesServer struct {
	grpc.ServerStream
}

func (x *deviceGatewayGetDevicesServer) Send(m *Device) error {
	return x.ServerStream.SendMsg(m)
}

func _DeviceGateway_GetDeviceByEndpoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceByEndpointRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceGatewayServer).GetDeviceByEndpoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.pb.DeviceGateway/GetDeviceByEndpoint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceGatewayServer).GetDeviceByEndpoint(ctx, req.(*GetDeviceByEndpointRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceGateway_GetResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceGatewayServer).GetResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.pb.DeviceGateway/GetResource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceGatewayServer).GetResource(ctx, req.(*GetResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceGateway_UpdateResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceGatewayServer).UpdateResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.pb.DeviceGateway/UpdateResource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceGatewayServer).UpdateResource(ctx, req.(*UpdateResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DeviceGateway_ServiceDesc is the grpc.ServiceDesc for DeviceGateway service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeviceGateway_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.pb.DeviceGateway",
	HandlerType: (*DeviceGatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDeviceByEndpoint",
			Handler:    _DeviceGateway_GetDeviceByEndpoint_Handler,
		},
		{
			MethodName: "GetResource",
			Handler:    _DeviceGateway_GetResource_Handler,
		},
		{
			MethodName: "UpdateResource",
			Handler:    _DeviceGateway_UpdateResource_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetDevices",
			Handler:       _DeviceGateway_GetDevices_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "github.com/plgd-dev/client-application/pb/service.proto",
}