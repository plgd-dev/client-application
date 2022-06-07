package pb

func (c *Content) Clone() *Content {
	if c == nil {
		return nil
	}
	return &Content{
		ContentType: c.GetContentType(),
		Data:        c.GetData(),
	}
}

func (d *Device) Clone() *Device {
	if d == nil {
		return nil
	}
	return &Device{
		Id:              d.GetId(),
		Types:           d.GetTypes(),
		Content:         d.GetContent().Clone(),
		OwnershipStatus: d.GetOwnershipStatus(),
		Endpoints:       d.GetEndpoints(),
	}
}
