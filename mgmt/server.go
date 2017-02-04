package mgmt

import (
	"fmt"
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/stats"
	"log"
	"net"
)

const (
	defaultAddr = "localhost"
	defaultPort = 8126
)

type Server struct {
	config  *config.Config
	handler *handler
}

func NewMgmtServer(
	appStats *stats.AppStats,
	processStats *stats.ProcessStats,
	config *config.Config,
	startupTime int64,
) *Server {
	return &Server{
		config: config,
		handler: &handler{
			appStats:     appStats,
			processStats: processStats,
			config:       config,
			startupTime:  startupTime,
		},
	}
}

func (s *Server) Start() {
	go s.start()
}

func (s *Server) start() {
	host := s.config.MgmtAddress
	if host == "" {
		host = defaultAddr
	}

	port := s.config.MgmtPort
	if port == 0 {
		port = defaultPort
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicln(err)
	}
	defer listener.Close()
	log.Printf("Management server listening on %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting", err.Error())
			continue
		}
		go s.handler.handleRequest(conn)
	}
}
