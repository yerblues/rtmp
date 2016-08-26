package rtmp

import (
	"net"
	"sync"
)

type Conn struct {
	conn     net.Conn
	isClient bool
	config   *Config

	handshakeL        sync.Mutex
	handshakeErr      error
	handshakeComplete bool
}

func (*Conn) ReadMessage() (Message, error) {
	return nil, nil
}

func (*Conn) WriteMessage(m Message) error {
	return nil
}

func (*Conn) Close() error {
	return nil
}

func (*Conn) Handshake() error {
	return nil
}
