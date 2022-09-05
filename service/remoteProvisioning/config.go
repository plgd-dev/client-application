package remoteProvisioning

import (
	"fmt"
	"time"

	"github.com/plgd-dev/client-application/pb"
)

type UserAgentConfig struct {
	CertificateAuthorityAddress string        `yaml:"certificateAuthorityAddress" json:"certificateAuthorityAddress"`
	CSRChallengeStateExpiration time.Duration `yaml:"csrChallengeStateExpiration" json:"csrChallengeStateExpiration"`
}

func (c *UserAgentConfig) ToProto() *pb.UserAgent {
	return &pb.UserAgent{
		CertificateAuthorityAddress: c.CertificateAuthorityAddress,
		CsrChallengeStateExpiration: c.CSRChallengeStateExpiration.Nanoseconds(),
	}
}

func (c *UserAgentConfig) Validate() error {
	if c.CertificateAuthorityAddress == "" {
		return fmt.Errorf("certificateAuthorityAddress('%v')", c.CertificateAuthorityAddress)
	}
	if c.CSRChallengeStateExpiration < 0 {
		return fmt.Errorf("csrChallengeStateExpiration('%v')", c.CSRChallengeStateExpiration)
	}
	return nil
}

type Mode string

func (m Mode) ToProto() pb.RemoteProvisioning_Mode {
	switch m {
	case Mode_None:
		return pb.RemoteProvisioning_MODE_NONE
	case Mode_UserAgent:
		return pb.RemoteProvisioning_USER_AGENT
	}
	return pb.RemoteProvisioning_MODE_NONE
}

const (
	Mode_None      Mode = "none"
	Mode_UserAgent Mode = "userAgent"
)

type AuthorizationConfig struct {
	Authority  string   `yaml:"authority" json:"authority"`
	ClientID   string   `yaml:"clientID" json:"clientId"`
	Audience   string   `yaml:"audience" json:"audience"`
	Scopes     []string `yaml:"scopes" json:"scopes"`
	OwnerClaim string   `yaml:"ownerClaim" json:"ownerClaim"`
}

func (m AuthorizationConfig) ToProto() *pb.Authorization {
	return &pb.Authorization{
		Authority:  m.Authority,
		OwnerClaim: m.OwnerClaim,
		ClientId:   m.ClientID,
		Audience:   m.Audience,
		Scopes:     m.Scopes,
	}
}

func (c *AuthorizationConfig) Validate(mode Mode) error {
	switch mode {
	case Mode_UserAgent:
		if c.ClientID == "" {
			return fmt.Errorf("clientID('%v')", c.ClientID)
		}
		if c.Authority == "" {
			return fmt.Errorf("authority('%v')", c.Authority)
		}
		if c.OwnerClaim == "" {
			return fmt.Errorf("ownerClaim('%v')", c.OwnerClaim)
		}
	case Mode_None:
	}
	return nil
}

type Config struct {
	Mode            Mode                `yaml:"mode" json:"mode"`
	UserAgentConfig UserAgentConfig     `yaml:"userAgent" json:"userAgent"`
	Authorization   AuthorizationConfig `yaml:"authorization" json:"authorization"`
}

func (c *Config) Validate() error {
	switch c.Mode {
	case Mode_None, "":
		c.Mode = Mode_None
	case Mode_UserAgent:
		if err := c.UserAgentConfig.Validate(); err != nil {
			return fmt.Errorf("userAgent.%w", err)
		}
		if err := c.Authorization.Validate(c.Mode); err != nil {
			return fmt.Errorf("authorization.%w", err)
		}
	}
	return nil
}

func (c *Config) ToProto() *pb.RemoteProvisioning {
	return &pb.RemoteProvisioning{
		Mode:          c.Mode.ToProto(),
		UserAgent:     c.UserAgentConfig.ToProto(),
		Authorization: c.Authorization.ToProto(),
	}
}
