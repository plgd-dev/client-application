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

package http

// HTTP Service URIs.
const (
	DeviceIDKey          = "deviceId"
	ResourceHrefKey      = "resourceHref"
	ResourcesPathKey     = "resources"
	ResourceLinksPathKey = "resource-links"

	ApplicationProtoJsonContentType = "application/protojson"
	ApplicationJsonContentType      = "application/json"

	Api       = "/api"
	ApiV1     = Api + "/v1"
	WellKnown = "/.well-known"
	Identity  = ApiV1 + "/identity"

	Devices             = ApiV1 + "/devices"
	Device              = Devices + "/{" + DeviceIDKey + "}"
	DeviceResourceLinks = Devices + "/{" + DeviceIDKey + "}/" + ResourceLinksPathKey
	DeviceResourceLink  = DeviceResourceLinks + "/{" + ResourceHrefKey + "}"
	DeviceResources     = Device + "/" + ResourcesPathKey
	DeviceResource      = DeviceResources + "/{" + ResourceHrefKey + "}"
	OwnDevice           = Device + "/own"
	DisownDevice        = Device + "/disown"

	Initialize             = ApiV1 + "/initialize"
	Reset                  = ApiV1 + "/reset"
	IdentityCertificate    = Identity + "/certificate"
	WellKnownJWKs          = WellKnown + "/jwks.json"
	WellKnownConfiguration = WellKnown + "/configuration"
)

func FinishInitialize(state string) string {
	return Initialize + "/" + state
}
