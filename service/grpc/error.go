// ************************************************************************
// Copyright (C) 2022 plgd.dev, s.r.o.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ************************************************************************

package grpc

import (
	"errors"

	coapStatus "github.com/plgd-dev/go-coap/v2/message/status"
	"github.com/plgd-dev/kit/v2/coapconv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convErrToGrpcStatus(defaultCode codes.Code, err error) *status.Status { //nolint:unparam
	var coapStatus coapStatus.Status
	if errors.As(err, &coapStatus) {
		return status.New(coapconv.ToGrpcCode(coapStatus.Code(), defaultCode), err.Error())
	}
	return status.New(defaultCode, err.Error())
}
