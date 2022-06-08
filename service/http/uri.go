package http

// HTTP Service URIs.
const (
	DeviceIDKey      = "deviceId"
	ResourceHrefKey  = "resourceHref"
	ResourcesPathKey = "resources"

	ApplicationProtoJsonContentType = "application/protojson"
	ApplicationJsonContentType      = "application/json"

	Api   = "/api"
	ApiV1 = Api + "/v1"

	Devices         = ApiV1 + "/devices"
	Device          = Devices + "/{" + DeviceIDKey + "}"
	DeviceResources = Device + "/" + ResourcesPathKey
	DeviceResource  = DeviceResources + "/{" + ResourceHrefKey + "}"
)
