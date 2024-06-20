package minecraft

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"time"
)

var proxies []proxyInfo

func init() {
	proxies, _ = loadProxies("proxies.txt")
}

func Connect(target string) (*minecraft.Conn, error) {
	if len(proxies) > 0 {
		// Override the default RakNet network with our anonymous RakNet network.
		minecraft.RegisterNetwork("raknet", NewAnonymousRakNet(proxies))
	}

	serverConn, err := minecraft.Dialer{
		TokenSource: src,
	}.DialTimeout("raknet", target, time.Minute*3)
	if err != nil {
		return nil, err
	}

	return serverConn, serverConn.DoSpawn()
}
