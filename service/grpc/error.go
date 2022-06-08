package grpc

import (
	"errors"

	coapStatus "github.com/plgd-dev/go-coap/v2/message/status"
	"github.com/plgd-dev/kit/v2/coapconv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convErrToGrpcStatus(defaultCode codes.Code, err error) *status.Status {
	var coapStatus coapStatus.Status
	if errors.As(err, &coapStatus) {
		return status.New(coapconv.ToGrpcCode(coapStatus.Code(), defaultCode), err.Error())
	}
	return status.New(defaultCode, err.Error())
}
