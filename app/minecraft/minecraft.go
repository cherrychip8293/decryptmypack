package minecraft

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"strings"
	"time"
)

func Connect(target string) (*minecraft.Conn, error) {
	if len(strings.Split(target, ":")) < 2 {
		target = target + ":19132"
	}

	serverConn, err := minecraft.Dialer{
		TokenSource: src,
	}.DialTimeout("raknet", target, time.Minute*3)
	if err != nil {
		return nil, err
	}

	return serverConn, serverConn.DoSpawn()
}
