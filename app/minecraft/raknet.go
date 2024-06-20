package minecraft

import (
	"context"
	"github.com/sandertv/go-raknet"
	"github.com/sandertv/gophertunnel/minecraft"
	"net"
)

// AnonymousRakNet ...
type AnonymousRakNet struct {
	minecraft.RakNet

	proxies []proxyInfo
}

// NewAnonymousRakNet ...
func NewAnonymousRakNet(proxies []proxyInfo) *AnonymousRakNet {
	return &AnonymousRakNet{proxies: proxies}
}

// DialContext ...
func (a *AnonymousRakNet) DialContext(ctx context.Context, address string) (net.Conn, error) {
	return raknet.Dialer{
		UpstreamDialer: &upDialer{a},
	}.DialContext(ctx, address)
}

type upDialer struct {
	*AnonymousRakNet
}

func (u *upDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return u.AnonymousRakNet.DialContext(ctx, address)
}

// Dial ...
func (a *AnonymousRakNet) Dial(network, address string) (net.Conn, error) {
	client := newClientWithProxy(randomProxy(a.proxies))
	return client.Dial(network, address)
}
