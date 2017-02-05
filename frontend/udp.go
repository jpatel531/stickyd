package frontend

import (
	"github.com/jpatel531/stickyd/config"
	"log"
	"net"
	"strconv"
)

const bufLen = 512

type udp struct{}

func (u *udp) Start(config *config.Frontend, handler Handler) {
	go u.start(config, handler)
}

func (u *udp) start(config *config.Frontend, handler Handler) {
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

	log.Println("UDP stats frontend listening on", udpAddr.String())
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

	go handler.HandleMessage(msg, addr)

	conn.WriteToUDP([]byte(""), addr)
}
