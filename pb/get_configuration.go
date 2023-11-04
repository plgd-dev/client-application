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
	"os"
	"time"

	pb "github.com/plgd-dev/hub/v2/grpc-gateway/pb"
	"github.com/plgd-dev/kit/v2/security"
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
		DeviceOauthClient:      c.GetDeviceOauthClient().Clone(),
		CertificateAuthority:   c.GetCertificateAuthority(),
	}
}

func (c *BuildInfo) Clone() *BuildInfo {
	if c == nil {
		return nil
	}
	return &BuildInfo{
		Version:    c.GetVersion(),
		BuildDate:  c.GetBuildDate(),
		CommitHash: c.GetCommitHash(),
		CommitDate: c.GetCommitDate(),
		ReleaseUrl: c.GetReleaseUrl(),
	}
}

func NewGetConfigurationResponse(info *BuildInfo) *GetConfigurationResponse {
	return &GetConfigurationResponse{
		BuildInfo:  info,
		Version:    info.GetVersion(),
		BuildDate:  info.GetBuildDate(),
		CommitHash: info.GetCommitHash(),
		CommitDate: info.GetCommitDate(),
		ReleaseUrl: info.GetReleaseUrl(),
	}
}

func (r *GetConfigurationResponse) Clone() *GetConfigurationResponse {
	if r == nil {
		return nil
	}
	return &GetConfigurationResponse{
		Version:            r.GetBuildInfo().GetVersion(),
		BuildDate:          r.GetBuildInfo().GetBuildDate(),
		CommitHash:         r.GetBuildInfo().GetCommitHash(),
		CommitDate:         r.GetBuildInfo().GetCommitDate(),
		ReleaseUrl:         r.GetBuildInfo().GetReleaseUrl(),
		RemoteProvisioning: r.RemoteProvisioning.Clone(),
		BuildInfo:          r.GetBuildInfo().Clone(),
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

func validateCA(path string) ([]byte, error) {
	data, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		return nil, err
	}
	_, err = security.ParseX509FromPEM(data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func validateCAPool(paths []string) ([]byte, error) {
	if len(paths) == 0 {
		return nil, nil
	}
	certificateAuthorities := make([]byte, 0, 512)
	for i, ca := range paths {
		data, err := validateCA(ca)
		if err != nil {
			return nil, fmt.Errorf("caPool[%v]('%v') - %w", i, paths, err)
		}
		certificateAuthorities = append(certificateAuthorities, data...)
		if certificateAuthorities[len(certificateAuthorities)-1] != '\n' {
			certificateAuthorities = append(certificateAuthorities, '\n')
		}
	}
	return certificateAuthorities, nil
}

func (c *RemoteProvisioning) validateForUserAgent() error {
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
	return nil
}

func (c *RemoteProvisioning) Validate() error {
	if c == nil {
		return nil
	}
	certificateAuthorities, err := validateCAPool(c.GetCaPool())
	if err != nil {
		return err
	}
	c.CertificateAuthorities = string(certificateAuthorities)
	switch c.GetMode() {
	case RemoteProvisioning_USER_AGENT:
		if err := c.validateForUserAgent(); err != nil {
			return err
		}
	case RemoteProvisioning_MODE_NONE:
	}
	return nil
}
