package frontends

import (
	"github.com/jpatel531/stickyd/config"
)

type Frontend interface {
	Start(*config.Frontend, Handler)
}

type RemoteInfo struct {
	Host string
	Port int
}

type Handler interface {
	HandleMessage(msg []byte, rinfo *RemoteInfo)
}
