// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	pb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	events "github.com/plgd-dev/hub/v2/resource-aggregate/events"
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
	// Get device information from the device. Device needs to be stored in cache otherwise it returns not found.
	GetDevice(ctx context.Context, in *GetDeviceRequest, opts ...grpc.CallOption) (*pb.Device, error)
	// Get resource links of devices. Device needs to be stored in cache otherwise it returns not found.
	GetDeviceResourceLinks(ctx context.Context, in *GetDeviceResourceLinksRequest, opts ...grpc.CallOption) (*events.ResourceLinksPublished, error)
	// Get resource from the device. Device needs to be stored in cache otherwise it returns not found.
	GetResource(ctx context.Context, in *GetResourceRequest, opts ...grpc.CallOption) (*pb.Resource, error)
	// Update resource at the device. Device needs to be stored in cache otherwise it returns not found.
	UpdateResource(ctx context.Context, in *pb.UpdateResourceRequest, opts ...grpc.CallOption) (*pb.UpdateResourceResponse, error)
	// Own the device. Device needs to be stored in cache otherwise it returns not found.
	OwnDevice(ctx context.Context, in *OwnDeviceRequest, opts ...grpc.CallOption) (*OwnDeviceResponse, error)
	// Disown the device. Device needs to be stored in cache otherwise it returns not found.
	DisownDevice(ctx context.Context, in *DisownDeviceRequest, opts ...grpc.CallOption) (*DisownDeviceResponse, error)
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
	Recv() (*pb.Device, error)
	grpc.ClientStream
}

type deviceGatewayGetDevicesClient struct {
	grpc.ClientStream
}

func (x *deviceGatewayGetDevicesClient) Recv() (*pb.Device, error) {
	m := new(pb.Device)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *deviceGatewayClient) GetDevice(ctx context.Context, in *GetDeviceRequest, opts ...grpc.CallOption) (*pb.Device, error) {
	out := new(pb.Device)
	err := c.cc.Invoke(ctx, "/service.pb.DeviceGateway/GetDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceGatewayClient) GetDeviceResourceLinks(ctx context.Context, in *GetDeviceResourceLinksRequest, opts ...grpc.CallOption) (*events.ResourceLinksPublished, error) {
	out := new(events.ResourceLinksPublished)
	err := c.cc.Invoke(ctx, "/service.pb.DeviceGateway/GetDeviceResourceLinks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceGatewayClient) GetResource(ctx context.Context, in *GetResourceRequest, opts ...grpc.CallOption) (*pb.Resource, error) {
	out := new(pb.Resource)
	err := c.cc.Invoke(ctx, "/service.pb.DeviceGateway/GetResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceGatewayClient) UpdateResource(ctx context.Context, in *pb.UpdateResourceRequest, opts ...grpc.CallOption) (*pb.UpdateResourceResponse, error) {
	out := new(pb.UpdateResourceResponse)
	err := c.cc.Invoke(ctx, "/service.pb.DeviceGateway/UpdateResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceGatewayClient) OwnDevice(ctx context.Context, in *OwnDeviceRequest, opts ...grpc.CallOption) (*OwnDeviceResponse, error) {
	out := new(OwnDeviceResponse)
	err := c.cc.Invoke(ctx, "/service.pb.DeviceGateway/OwnDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceGatewayClient) DisownDevice(ctx context.Context, in *DisownDeviceRequest, opts ...grpc.CallOption) (*DisownDeviceResponse, error) {
	out := new(DisownDeviceResponse)
	err := c.cc.Invoke(ctx, "/service.pb.DeviceGateway/DisownDevice", in, out, opts...)
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
	// Get device information from the device. Device needs to be stored in cache otherwise it returns not found.
	GetDevice(context.Context, *GetDeviceRequest) (*pb.Device, error)
	// Get resource links of devices. Device needs to be stored in cache otherwise it returns not found.
	GetDeviceResourceLinks(context.Context, *GetDeviceResourceLinksRequest) (*events.ResourceLinksPublished, error)
	// Get resource from the device. Device needs to be stored in cache otherwise it returns not found.
	GetResource(context.Context, *GetResourceRequest) (*pb.Resource, error)
	// Update resource at the device. Device needs to be stored in cache otherwise it returns not found.
	UpdateResource(context.Context, *pb.UpdateResourceRequest) (*pb.UpdateResourceResponse, error)
	// Own the device. Device needs to be stored in cache otherwise it returns not found.
	OwnDevice(context.Context, *OwnDeviceRequest) (*OwnDeviceResponse, error)
	// Disown the device. Device needs to be stored in cache otherwise it returns not found.
	DisownDevice(context.Context, *DisownDeviceRequest) (*DisownDeviceResponse, error)
	mustEmbedUnimplementedDeviceGatewayServer()
}

// UnimplementedDeviceGatewayServer must be embedded to have forward compatible implementations.
type UnimplementedDeviceGatewayServer struct {
}

func (UnimplementedDeviceGatewayServer) GetDevices(*GetDevicesRequest, DeviceGateway_GetDevicesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetDevices not implemented")
}
func (UnimplementedDeviceGatewayServer) GetDevice(context.Context, *GetDeviceRequest) (*pb.Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDevice not implemented")
}
func (UnimplementedDeviceGatewayServer) GetDeviceResourceLinks(context.Context, *GetDeviceResourceLinksRequest) (*events.ResourceLinksPublished, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeviceResourceLinks not implemented")
}
func (UnimplementedDeviceGatewayServer) GetResource(context.Context, *GetResourceRequest) (*pb.Resource, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResource not implemented")
}
func (UnimplementedDeviceGatewayServer) UpdateResource(context.Context, *pb.UpdateResourceRequest) (*pb.UpdateResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateResource not implemented")
}
func (UnimplementedDeviceGatewayServer) OwnDevice(context.Context, *OwnDeviceRequest) (*OwnDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OwnDevice not implemented")
}
func (UnimplementedDeviceGatewayServer) DisownDevice(context.Context, *DisownDeviceRequest) (*DisownDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisownDevice not implemented")
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
	Send(*pb.Device) error
	grpc.ServerStream
}

type deviceGatewayGetDevicesServer struct {
	grpc.ServerStream
}

func (x *deviceGatewayGetDevicesServer) Send(m *pb.Device) error {
	return x.ServerStream.SendMsg(m)
}

func _DeviceGateway_GetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceGatewayServer).GetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.pb.DeviceGateway/GetDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceGatewayServer).GetDevice(ctx, req.(*GetDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceGateway_GetDeviceResourceLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceResourceLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceGatewayServer).GetDeviceResourceLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.pb.DeviceGateway/GetDeviceResourceLinks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceGatewayServer).GetDeviceResourceLinks(ctx, req.(*GetDeviceResourceLinksRequest))
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
	in := new(pb.UpdateResourceRequest)
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
		return srv.(DeviceGatewayServer).UpdateResource(ctx, req.(*pb.UpdateResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceGateway_OwnDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OwnDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceGatewayServer).OwnDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.pb.DeviceGateway/OwnDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceGatewayServer).OwnDevice(ctx, req.(*OwnDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceGateway_DisownDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisownDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceGatewayServer).DisownDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.pb.DeviceGateway/DisownDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceGatewayServer).DisownDevice(ctx, req.(*DisownDeviceRequest))
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
			MethodName: "GetDevice",
			Handler:    _DeviceGateway_GetDevice_Handler,
		},
		{
			MethodName: "GetDeviceResourceLinks",
			Handler:    _DeviceGateway_GetDeviceResourceLinks_Handler,
		},
		{
			MethodName: "GetResource",
			Handler:    _DeviceGateway_GetResource_Handler,
		},
		{
			MethodName: "UpdateResource",
			Handler:    _DeviceGateway_UpdateResource_Handler,
		},
		{
			MethodName: "OwnDevice",
			Handler:    _DeviceGateway_OwnDevice_Handler,
		},
		{
			MethodName: "DisownDevice",
			Handler:    _DeviceGateway_DisownDevice_Handler,
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
