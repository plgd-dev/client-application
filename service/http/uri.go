package http

// HTTP Service URIs.
const (
	DeviceIDKey          = "deviceId"
	ResourceHrefKey      = "resourceHref"
	ResourcesPathKey     = "resources"
	ResourceLinksPathKey = "resource-links"

	ApplicationProtoJsonContentType = "application/protojson"
	ApplicationJsonContentType      = "application/json"

	Api   = "/api"
	ApiV1 = Api + "/v1"

	Devices             = ApiV1 + "/devices"
	Device              = Devices + "/{" + DeviceIDKey + "}"
	DeviceResourceLinks = Devices + "/{" + DeviceIDKey + "}/" + ResourceLinksPathKey
	DeviceResources     = Device + "/" + ResourcesPathKey
	DeviceResource      = DeviceResources + "/{" + ResourceHrefKey + "}"
	OwnDevice           = Device + "/own"
	DisownDevice        = Device + "/disown"
)
