//go:build !telio

package main

import (
	"errors"
	"net/netip"

	"github.com/NordSecurity/nordvpn-linux/config"
	"github.com/NordSecurity/nordvpn-linux/config/remote"
	cesh "github.com/NordSecurity/nordvpn-linux/core/mesh"
	"github.com/NordSecurity/nordvpn-linux/daemon"
	"github.com/NordSecurity/nordvpn-linux/daemon/vpn"
	"github.com/NordSecurity/nordvpn-linux/daemon/vpn/nordlynx"
	"github.com/NordSecurity/nordvpn-linux/daemon/vpn/openvpn"
	"github.com/NordSecurity/nordvpn-linux/meshnet"
	"github.com/NordSecurity/nordvpn-linux/tunnel"
)

func getVpnFactory(eventsDbPath string, fwmark uint32, envIsDev bool,
	telioCfg remote.RemoteConfigGetter, deviceID, appVersion string) daemon.FactoryFunc {
	return func(tech config.Technology) (vpn.VPN, error) {
		switch tech {
		case config.Technology_NORDLYNX:
			return nordlynx.NewKernelSpace(fwmark), nil
		case config.Technology_OPENVPN:
			return openvpn.New(fwmark), nil
		case config.Technology_UNKNOWN_TECHNOLOGY:
			fallthrough
		default:
			return nil, errors.New("no such technology")
		}
	}
}

// noopMesh is a noop implementation of meshnet. It is used when telio
// is not available and should be used only for development purposes
type noopMesh bool

func (noopMesh) Enable(netip.Addr, string) error { return nil }
func (noopMesh) Disable() error                  { return nil }
func (noopMesh) IsActive() bool                  { return false }
func (noopMesh) Refresh(cesh.MachineMap) error   { return nil }
func (noopMesh) Tun() tunnel.T                   { return &tunnel.Tunnel{} }
func (noopMesh) StatusMap() (map[string]string, error) {
	return map[string]string{}, nil
}

func meshnetImplementation(fn daemon.FactoryFunc) (meshnet.Mesh, error) {
	return noopMesh(true), nil
}

type noopKeygen bool

func (noopKeygen) Private() string      { return "" }
func (noopKeygen) Public(string) string { return "" }

func keygenImplementation(daemon.FactoryFunc) (noopKeygen, error) { return noopKeygen(true), nil }
