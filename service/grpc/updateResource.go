package grpc

import (
	"context"
	"fmt"

	"github.com/plgd-dev/client-application/pkg/rawcodec"
	"github.com/plgd-dev/go-coap/v2/message"
	grpcgwPb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/hub/v2/resource-aggregate/commands"
	"github.com/plgd-dev/hub/v2/resource-aggregate/events"
	"github.com/plgd-dev/kit/v2/codec/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convContentToOcfCbor(content *grpcgwPb.Content) ([]byte, error) {
	switch content.GetContentType() {
	case message.AppCBOR.String(), message.AppOcfCbor.String():
		return content.GetData(), nil
	case message.AppJSON.String():
		data, err := json.ToCBOR(string(content.GetData()))
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "cannot convert json to cbor: %v", err)
		}
		return data, nil
	}
	return nil, status.Errorf(codes.InvalidArgument, "unsupported content type '%v'", content.GetContentType())
}

func (s *DeviceGatewayServer) UpdateResource(ctx context.Context, req *grpcgwPb.UpdateResourceRequest) (*grpcgwPb.UpdateResourceResponse, error) {
	updateData, err := convContentToOcfCbor(req.GetContent())
	if err != nil {
		return nil, err
	}
	dev, err := s.getDevice(req.GetResourceId().GetDeviceId())
	if err != nil {
		return nil, err
	}
	link, err := dev.getResourceLink(ctx, req.GetResourceId())
	if err != nil {
		return nil, err
	}
	if dev.ToProto().OwnershipStatus != grpcgwPb.Device_OWNED && len(link.Endpoints.FilterUnsecureEndpoints()) == 0 {
		return nil, status.Error(codes.PermissionDenied, "device is not owned")
	}
	codec := rawcodec.GetRawCodec(message.AppOcfCbor)
	var data []byte
	err = dev.UpdateResourceWithCodec(ctx, link, codec, updateData, &data)
	if err != nil {
		return nil, convErrToGrpcStatus(codes.Unavailable, fmt.Errorf("cannot get resource %v for device %v: %w", req.GetResourceId().GetHref(), dev.ID, err)).Err()
	}
	contentType := ""
	if len(data) > 0 {
		contentType = message.AppOcfCbor.String()
	}
	return &grpcgwPb.UpdateResourceResponse{
		Data: &events.ResourceUpdated{
			Content: &commands.Content{
				ContentType: contentType,
				Data:        data,
			},
			Status: commands.Status_OK,
		},
	}, nil
}