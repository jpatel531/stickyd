package frontends

import (
	"github.com/jpatel531/stickyd/config"
	"log"
	"net"
	"strconv"
)

const bufLen = 512

type UDP struct{}

func (u *UDP) Start(config *config.Frontend, handler Handler) {
	go u.start(config, handler)
}

func (u *UDP) start(config *config.Frontend, handler Handler) {
	var udpVersion string
	if config.AddressIPV6 {
		udpVersion = "udp6"
	} else {
		udpVersion = "udp4"
	}

	udpAddr, err := net.ResolveUDPAddr(udpVersion, ":"+strconv.Itoa(config.Port))
	if err != nil {
		log.Panicln(err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println(err)
	}

	for {
		readUDPMessages(conn, handler)
	}
}

func readUDPMessages(conn *net.UDPConn, handler Handler) {
	buf := make([]byte, bufLen)

	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Println(err)
		return
	}

	msg := buf[:n]

	go handler.HandleMessage(msg, &RemoteInfo{
		Host: addr.IP.String(),
		Port: addr.Port,
	})

	conn.WriteToUDP([]byte(""), addr)
}
