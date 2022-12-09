// ************************************************************************
// Copyright (C) 2022 plgd.dev, s.r.o.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ************************************************************************
package pb

import (
	"fmt"
	"time"

	pb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"gopkg.in/yaml.v3"
)

type (
	WebOauthClient    = pb.WebOAuthClient
	DeviceOauthClient = pb.DeviceOAuthClient
)

func (c *UserAgent) Clone() *UserAgent {
	if c == nil {
		return nil
	}
	return &UserAgent{
		CsrChallengeStateExpiration: c.CsrChallengeStateExpiration,
	}
}

func (c *RemoteProvisioning) Clone() *RemoteProvisioning {
	if c == nil {
		return nil
	}
	return &RemoteProvisioning{
		CurrentTime:            c.GetCurrentTime(),
		Mode:                   c.GetMode(),
		UserAgent:              c.GetUserAgent().Clone(),
		WebOauthClient:         c.GetWebOauthClient().Clone(),
		JwtOwnerClaim:          c.GetJwtOwnerClaim(),
		Id:                     c.GetId(),
		CoapGateway:            c.GetCoapGateway(),
		CertificateAuthorities: c.GetCertificateAuthorities(),
		Authority:              c.GetAuthority(),
		HttpGatewayAddress:     c.GetHttpGatewayAddress(),
		DeviceOauthClient:      c.GetDeviceOauthClient().Clone(),
		CertificateAuthority:   c.GetCertificateAuthority(),
	}
}

func (r *GetConfigurationResponse) Clone() *GetConfigurationResponse {
	if r == nil {
		return nil
	}
	return &GetConfigurationResponse{
		Version:            r.Version,
		BuildDate:          r.BuildDate,
		CommitHash:         r.CommitHash,
		CommitDate:         r.CommitDate,
		ReleaseUrl:         r.ReleaseUrl,
		RemoteProvisioning: r.RemoteProvisioning.Clone(),
	}
}

func ValidateWebOAuthClient(c *pb.WebOAuthClient) error {
	if c.GetClientId() == "" {
		return fmt.Errorf("clientID('%v')", c.GetClientId())
	}
	return nil
}

func (c *UserAgent) Validate() error {
	if c.GetCsrChallengeStateExpiration() == 0 {
		return fmt.Errorf("csrChallengeStateExpiration('%v')", c.GetCsrChallengeStateExpiration())
	}
	return nil
}

func (c *UserAgent) UnmarshalYAML(value *yaml.Node) error {
	var v struct {
		CsrChallengeStateExpiration time.Duration `yaml:"csrChallengeStateExpiration"`
	}
	if err := value.Decode(&v); err != nil {
		return fmt.Errorf("csrChallengeStateExpiration('%v') - %w", c.GetCsrChallengeStateExpiration(), err)
	}
	c.CsrChallengeStateExpiration = v.CsrChallengeStateExpiration.Nanoseconds()
	return nil
}

func (c *UserAgent) MarshalYAML() (interface{}, error) {
	var v struct {
		CsrChallengeStateExpiration time.Duration `yaml:"csrChallengeStateExpiration"`
	}
	v.CsrChallengeStateExpiration = time.Nanosecond * time.Duration(c.GetCsrChallengeStateExpiration())
	return v, nil
}

func (c RemoteProvisioning_Mode) MarshalYAML() (interface{}, error) {
	switch c {
	case RemoteProvisioning_USER_AGENT:
		return "userAgent", nil
	case RemoteProvisioning_MODE_NONE:
		return "", nil
	}
	return "", nil
}

func (c *RemoteProvisioning_Mode) UnmarshalYAML(value *yaml.Node) error {
	var v string
	if err := value.Decode(&v); err != nil {
		return err
	}
	if v == "userAgent" {
		*c = RemoteProvisioning_USER_AGENT
		return nil
	}
	*c = RemoteProvisioning_MODE_NONE
	return nil
}

func (c *RemoteProvisioning) Validate() error {
	if c == nil {
		return nil
	}
	switch c.GetMode() {
	case RemoteProvisioning_USER_AGENT:
		if err := ValidateWebOAuthClient(c.GetWebOauthClient()); err != nil {
			return fmt.Errorf("webOAuthClient.%w", err)
		}
		if c.GetAuthority() == "" {
			return fmt.Errorf("authority('%v')", c.GetAuthority())
		}
		if c.GetJwtOwnerClaim() == "" {
			return fmt.Errorf("ownerClaim('%v')", c.GetJwtOwnerClaim())
		}
		if err := c.GetUserAgent().Validate(); err != nil {
			return fmt.Errorf("userAgent.%w", err)
		}
		if c.GetCertificateAuthorities() == "" {
			return fmt.Errorf("certificateAuthorities('%v')", c.GetCertificateAuthorities())
		}
	case RemoteProvisioning_MODE_NONE:
	}
	return nil
}
