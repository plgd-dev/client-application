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

package rawcodec

import (
	kitNetCoap "github.com/plgd-dev/device/pkg/net/coap"
	"github.com/plgd-dev/go-coap/v2/message"
	codecOcf "github.com/plgd-dev/kit/v2/codec/ocf"
)

// GetRawCodec returns raw codec depends on contentFormat.
func GetRawCodec(contentFormat message.MediaType) kitNetCoap.Codec {
	if contentFormat == message.AppCBOR || contentFormat == message.AppOcfCbor {
		return codecOcf.RawVNDOCFCBORCodec{}
	}
	return codecOcf.NoCodec{
		MediaType: uint16(contentFormat),
	}
}
