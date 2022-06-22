package pb

import (
	"fmt"
	"net"
	"net/netip"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GetDevicesRequestConfig struct {
	// Devices are taken from the cache. Default: false.
	UseCache bool `yaml:"useCache"`
	// Filter by multicast IP address version. Default: [] - multicast is disabled. If it is set, the new devices will be added to cache.
	UseMulticast []string `yaml:"useMulticast"`

	// Returns devices via endpoints. Default: [] - filter is disabled. New devices will be added to cache. Not reachable devices will be not in response.
	// Endpoint can be in format:
	// - <host>:<port> is interpreted as coap://<host>:<port>
	// - <host> is interpreted as coap://<host>:5683
	UseEndpoints []string `yaml:"useEndpoints"`

	// How long to wait for the devices responses for responses in milliseconds. Default: 0 - means 1sec.
	Timeout time.Duration `yaml:"timeout"`

	// Filter by ownership status. Default: [UNOWNED,OWNED].
	OwnershipStatusFilter []string `yaml:"ownershipStatusFilter"`

	// Filter by device resource type of oic/d. Default: [] - filter is disabled.
	TypeFilter []string `yaml:"typeFilter"`
}

func validateHostPort(domain *regexp.Regexp, host, portStr string) error {
	if host == "" {
		return fmt.Errorf("empty host")
	}
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return fmt.Errorf("invalid port")
	}
	if port == 0 {
		return fmt.Errorf("invalid port")
	}
	if addr, err := netip.ParseAddr(host); err == nil && !addr.IsUnspecified() && !addr.IsMulticast() {
		return nil
	}
	matched := domain.MatchString(host)
	if !matched {
		return fmt.Errorf("invalid domain")
	}
	return nil
}

func validateEndpoint(domain *regexp.Regexp, value string) error {
	if addr, err := netip.ParseAddrPort(value); err == nil && !addr.Addr().IsUnspecified() && !addr.Addr().IsMulticast() && addr.Port() > 0 {
		return nil
	}
	host, portStr, err := net.SplitHostPort(value)
	if err != nil && strings.Contains(err.Error(), "missing port in address") {
		host = value
		portStr = "1234"
		err = nil
	}
	if err != nil {
		return err
	}
	if len(host) >= 2 && host[0] == '[' && host[len(host)-1] == ']' {
		host = host[1 : len(host)-1]
	}
	return validateHostPort(domain, host, portStr)
}

func (c *GetDevicesRequestConfig) Validate() error {
	reDomain := regexp.MustCompile(`([A-Za-z0-9][A-Za-z0-9\-]{1,61}[A-Za-z0-9]\.{0,1})+`)
	for _, value := range c.UseEndpoints {
		err := validateEndpoint(reDomain, value)
		if err != nil {
			return fmt.Errorf("useEndpoints('%v') - %w", value, err)
		}
	}
	for _, value := range c.UseMulticast {
		_, ok := GetDevicesRequest_UseMulticast_value[strings.ToUpper(value)]
		if !ok {
			return fmt.Errorf("useMulticast('%v')", value)
		}
	}
	for _, value := range c.OwnershipStatusFilter {
		_, ok := GetDevicesRequest_OwnershipStatusFilter_value[strings.ToUpper(value)]
		if !ok {
			return fmt.Errorf("ownershipStatusFilter('%v')", value)
		}
	}
	if c.Timeout < time.Millisecond*50 {
		return fmt.Errorf("timeout('%v')", c.Timeout)
	}
	return nil
}

func (c *GetDevicesRequestConfig) ToGetDevicesRequest() *GetDevicesRequest {
	r := GetDevicesRequest{
		UseCache:     c.UseCache,
		TypeFilter:   c.TypeFilter,
		UseEndpoints: c.UseEndpoints,
		Timeout:      uint32(c.Timeout.Milliseconds()),
	}
	for _, value := range c.UseMulticast {
		r.UseMulticast = append(r.UseMulticast, GetDevicesRequest_UseMulticast(GetDevicesRequest_UseMulticast_value[strings.ToUpper(value)]))
	}
	for _, value := range c.OwnershipStatusFilter {
		r.OwnershipStatusFilter = append(r.OwnershipStatusFilter, GetDevicesRequest_OwnershipStatusFilter(GetDevicesRequest_OwnershipStatusFilter_value[strings.ToUpper(value)]))
	}
	return &r
}
