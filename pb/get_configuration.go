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

func (r *UserAgent) Clone() *UserAgent {
	if r == nil {
		return nil
	}
	return &UserAgent{
		CertificateAuthorityAddress: r.CertificateAuthorityAddress,
		CsrChallengeStateExpiration: r.CsrChallengeStateExpiration,
	}
}

func (r *Authorization) Clone() *Authorization {
	if r == nil {
		return nil
	}
	scopes := make([]string, len(r.Scopes))
	copy(scopes, r.Scopes)
	return &Authorization{
		ClientId:   r.ClientId,
		Audience:   r.Audience,
		Scopes:     scopes,
		OwnerClaim: r.OwnerClaim,
		Authority:  r.Authority,
	}
}

func (r *RemoteProvisioning) Clone() *RemoteProvisioning {
	if r == nil {
		return nil
	}
	return &RemoteProvisioning{
		Mode:          r.Mode,
		UserAgent:     r.UserAgent.Clone(),
		Authorization: r.Authorization.Clone(),
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
