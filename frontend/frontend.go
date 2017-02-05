package frontend

import (
	"github.com/jpatel531/stickyd/config"
	"net"
)

type Frontend interface {
	Start(*config.Frontend, Handler)
}

type Handler interface {
	HandleMessage(msg []byte, addr net.Addr)
}

var Frontends = map[string]Frontend{
	"udp": &udp{},
	"tcp": &tcp{},
}
