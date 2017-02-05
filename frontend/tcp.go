package frontend

import (
	"fmt"
	"github.com/jpatel531/stickyd/config"
	"log"
	"net"
)

const (
	defaultAddr = "localhost"
)

type tcp struct{}

func (t *tcp) Start(cfg *config.Frontend, handler Handler) {
	go t.start(cfg, handler)
}

func (t *tcp) start(cfg *config.Frontend, handler Handler) {
	host := cfg.Address
	if host == "" {
		host = defaultAddr
	}

	port := cfg.Port
	if port == 0 {
		log.Panicln("[tcp frontend] unspecified port")
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicln(err)
	}
	defer listener.Close()

	log.Println("TCP stats frontend listening on", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting", err.Error())
			continue
		}

		handleConnection(conn, handler)
	}
}

func handleConnection(conn net.Conn, handler Handler) {
	defer conn.Close()

	buf := make([]byte, bufLen)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading:", err.Error())
		return
	}
	msg := buf[:n]

	go handler.HandleMessage(msg, conn.RemoteAddr())
}
