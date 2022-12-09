package remoteProvisioning

import (
	"time"

	"github.com/plgd-dev/client-application/pb"
)

type Config = pb.RemoteProvisioning

var defaultConfig = Config{
	Mode: pb.RemoteProvisioning_MODE_NONE,
	UserAgent: &pb.UserAgent{
		CsrChallengeStateExpiration: time.Minute.Nanoseconds(),
	},
	JwtOwnerClaim:  "sub",
	WebOauthClient: &pb.WebOauthClient{
		//	Scopes: []string{},
	},
	DeviceOauthClient: &pb.DeviceOauthClient{
		//	Scopes: []string{},
	},
}

func DefaultConfig() *Config {
	return defaultConfig.Clone()
}
