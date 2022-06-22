package pb_test

import (
	"testing"
	"time"

	"github.com/plgd-dev/client-application/pb"
	"github.com/stretchr/testify/require"
)

func TestGetDevicesRequestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		value   pb.GetDevicesRequestConfig
		wantErr bool
	}{
		{
			name: "ok",
			value: pb.GetDevicesRequestConfig{
				UseMulticast:          []string{"ipv4", "ipv6"},
				UseEndpoints:          []string{"plgd", "plgd.dev", "plgd.dev:1234", "192.168.1.1", "192.168.1.1:1234", "[::1]", "[::1]:1234"},
				OwnershipStatusFilter: []string{"unowned", "owned"},
				Timeout:               time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.value.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
