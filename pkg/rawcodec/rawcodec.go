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
