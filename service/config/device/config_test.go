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

package device_test

import (
	"testing"

	"github.com/plgd-dev/client-application/service/config/device"
	"github.com/plgd-dev/client-application/test"
	"github.com/stretchr/testify/require"
)

func TestManufacturerTLSConfigValidate(t *testing.T) {
	type fields struct {
		CAPool   string
		KeyFile  string
		CertFile string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				CAPool:   test.MFG_ROOT_CA_CRT,
				KeyFile:  test.MFG_CLIENT_APPLICATION_KEY,
				CertFile: test.MFG_CLIENT_APPLICATION_CRT,
			},
		},
		{
			name: "invalid caPool",
			fields: fields{
				CAPool:   "",
				KeyFile:  test.MFG_CLIENT_APPLICATION_KEY,
				CertFile: test.MFG_CLIENT_APPLICATION_CRT,
			},
			wantErr: true,
		},
		{
			name: "invalid certFile",
			fields: fields{
				CAPool:   test.MFG_ROOT_CA_CRT,
				KeyFile:  test.MFG_CLIENT_APPLICATION_KEY,
				CertFile: "",
			},
			wantErr: true,
		},
		{
			name: "invalid keyFile",
			fields: fields{
				CAPool:   test.MFG_ROOT_CA_CRT,
				KeyFile:  "",
				CertFile: test.MFG_CLIENT_APPLICATION_CRT,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &device.ManufacturerTLSConfig{
				CAPool:   tt.fields.CAPool,
				KeyFile:  tt.fields.KeyFile,
				CertFile: tt.fields.CertFile,
			}
			err := c.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestManufacturerConfigValidate(t *testing.T) {
	type fields struct {
		TLS device.ManufacturerTLSConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				TLS: device.ManufacturerTLSConfig{
					CAPool:   test.MFG_ROOT_CA_CRT,
					KeyFile:  test.MFG_CLIENT_APPLICATION_KEY,
					CertFile: test.MFG_CLIENT_APPLICATION_CRT,
				},
			},
		},
		{
			name:    "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &device.ManufacturerConfig{
				TLS: tt.fields.TLS,
			}
			err := c.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestOwnershipTransferConfigValidate(t *testing.T) {
	type fields struct {
		Methods      []device.OwnershipTransferMethod
		Manufacturer device.ManufacturerConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: string(device.OwnershipTransferJustWorks),
			fields: fields{
				Methods: []device.OwnershipTransferMethod{device.OwnershipTransferJustWorks},
			},
		},
		{
			name: string(device.OwnershipTransferManufacturerCertificate),
			fields: fields{
				Methods:      []device.OwnershipTransferMethod{device.OwnershipTransferManufacturerCertificate},
				Manufacturer: test.MakeDeviceConfig().COAP.OwnershipTransfer.Manufacturer,
			},
		},
		{
			name: "invalid method",
			fields: fields{
				Methods: []device.OwnershipTransferMethod{"invalid"},
			},
			wantErr: true,
		},
		{
			name: "invalid manufacturerCertificate",
			fields: fields{
				Methods: []device.OwnershipTransferMethod{device.OwnershipTransferManufacturerCertificate},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &device.OwnershipTransferConfig{
				Methods:      tt.fields.Methods,
				Manufacturer: tt.fields.Manufacturer,
			}
			err := c.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestCoapConfigValidate(t *testing.T) {
	type fields struct {
		MaxMessageSize    uint32
		InactivityMonitor device.InactivityMonitor
		BlockwiseTransfer device.BlockwiseTransferConfig
		OwnershipTransfer device.OwnershipTransferConfig
		TLS               device.TLSConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				MaxMessageSize:    test.MakeDeviceConfig().COAP.MaxMessageSize,
				InactivityMonitor: test.MakeDeviceConfig().COAP.InactivityMonitor,
				BlockwiseTransfer: test.MakeDeviceConfig().COAP.BlockwiseTransfer,
				OwnershipTransfer: test.MakeDeviceConfig().COAP.OwnershipTransfer,
				TLS:               test.MakeDeviceConfig().COAP.TLS,
			},
		},
		{
			name: "invalid inactivityMonitor",
			fields: fields{
				MaxMessageSize:    test.MakeDeviceConfig().COAP.MaxMessageSize,
				BlockwiseTransfer: test.MakeDeviceConfig().COAP.BlockwiseTransfer,
				OwnershipTransfer: test.MakeDeviceConfig().COAP.OwnershipTransfer,
				TLS:               test.MakeDeviceConfig().COAP.TLS,
			},
			wantErr: true,
		},
		{
			name: "invalid maxMessageSize",
			fields: fields{
				InactivityMonitor: test.MakeDeviceConfig().COAP.InactivityMonitor,
				BlockwiseTransfer: test.MakeDeviceConfig().COAP.BlockwiseTransfer,
				OwnershipTransfer: test.MakeDeviceConfig().COAP.OwnershipTransfer,
				TLS:               test.MakeDeviceConfig().COAP.TLS,
			},
			wantErr: true,
		},
		{
			name: "invalid blockwiseTransfer",
			fields: fields{
				MaxMessageSize:    test.MakeDeviceConfig().COAP.MaxMessageSize,
				InactivityMonitor: test.MakeDeviceConfig().COAP.InactivityMonitor,
				OwnershipTransfer: test.MakeDeviceConfig().COAP.OwnershipTransfer,
				TLS:               test.MakeDeviceConfig().COAP.TLS,
				BlockwiseTransfer: device.BlockwiseTransferConfig{
					Enabled: true,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid ownershipTransfer",
			fields: fields{
				MaxMessageSize:    test.MakeDeviceConfig().COAP.MaxMessageSize,
				InactivityMonitor: test.MakeDeviceConfig().COAP.InactivityMonitor,
				BlockwiseTransfer: test.MakeDeviceConfig().COAP.BlockwiseTransfer,
				TLS:               test.MakeDeviceConfig().COAP.TLS,
			},
			wantErr: true,
		},
		{
			name: "invalid tls",
			fields: fields{
				MaxMessageSize:    test.MakeDeviceConfig().COAP.MaxMessageSize,
				InactivityMonitor: test.MakeDeviceConfig().COAP.InactivityMonitor,
				BlockwiseTransfer: test.MakeDeviceConfig().COAP.BlockwiseTransfer,
				OwnershipTransfer: test.MakeDeviceConfig().COAP.OwnershipTransfer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &device.CoapConfig{
				MaxMessageSize:    tt.fields.MaxMessageSize,
				InactivityMonitor: tt.fields.InactivityMonitor,
				BlockwiseTransfer: tt.fields.BlockwiseTransfer,
				OwnershipTransfer: tt.fields.OwnershipTransfer,
				TLS:               tt.fields.TLS,
			}
			err := c.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
