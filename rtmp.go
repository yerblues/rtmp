package rtmp

import (
	"net"
)

func Server(conn net.Conn, config *Config) *Conn {
	return nil
}

func Client(conn net.Conn, config *Config) *Conn {
	return nil
}

func NewListener(inner net.Listener, config *Config) net.Listener {
	return nil
}

func Listen(network, laddr string, config *Config) (net.Listener, error) {
	return nil, nil
}

func DialWithDialer(dialer *net.Dialer, network, addr string, config *Config) (*Conn, error) {
	return nil, nil
}

func Dial(network, addr string, config *Config) (*Conn, error) {
	return nil, nil
}
