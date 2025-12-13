package server

import (
	"net"

	"github.com/rs/zerolog/log"
)

type TCPHandlerFunc func(net.Conn)

func (f TCPHandlerFunc) HandleTCP(conn net.Conn) {
	f(conn)
}

type TCPHandler interface {
	HandleTCP(net.Conn)
}

func TCP(h TCPHandler) {
	addr := GetListenAddr()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to listen")
	}
	log.Info().Stringer("addr", ln.Addr()).Msg("listening")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Err(err).Msg("accept err")
			continue
		}
		go h.HandleTCP(conn)
	}
}
