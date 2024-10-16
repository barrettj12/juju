// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing

type MockProxier struct {
	// See Proxier interface
	StartFn func() error

	// See Proxier interface
	StopFn func()

	// See Proxier interface
	TypeFn func() string

	// See Proxier interface
	MarshalYAMLFn func() (interface{}, error)

	// See Proxier interface
	RawConfigFn func() (map[string]interface{}, error)
}

type MockTunnelProxier struct {
	*MockProxier

	// See TunnelProxier interface
	HostFn func() string

	// See TunnelProxier interface
	PortFn func() string
}

func NewMockTunnelProxier() *MockTunnelProxier {
	return &MockTunnelProxier{
		MockProxier: &MockProxier{},
	}
}

func (mp *MockProxier) Start() error {
	if mp.StartFn == nil {
		return nil
	}
	return mp.StartFn()
}

func (mp *MockProxier) MarshalYAML() (interface{}, error) {
	if mp.MarshalYAMLFn == nil {
		return nil, nil
	}
	return mp.MarshalYAMLFn()
}

func (mp *MockProxier) Insecure() {}

func (mp *MockProxier) Stop() {
	if mp.StopFn != nil {
		mp.StopFn()
	}
}

func (mp *MockProxier) RawConfig() (map[string]interface{}, error) {
	if mp.RawConfigFn == nil {
		return nil, nil
	}
	return mp.RawConfigFn()
}

func (mp *MockProxier) Type() string {
	if mp.TypeFn == nil {
		return "mock-proxier"
	}
	return mp.TypeFn()
}

func (mtp *MockTunnelProxier) Host() string {
	if mtp.HostFn == nil {
		return ""
	}
	return mtp.HostFn()
}

func (mtp *MockTunnelProxier) Port() string {
	if mtp.PortFn == nil {
		return ""
	}
	return mtp.PortFn()
}
