package frontends

type Config struct {
	AddressIPV6 bool
	Address     string
	Port        int
}

type Frontend interface {
	Start(*Config, Handler)
}

type RemoteInfo struct {
	Host string
	Port int
}

type Handler interface {
	HandleMessage(msg []byte, rinfo *RemoteInfo)
}
