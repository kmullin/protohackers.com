package server

import (
	"net"
	"os"

	"github.com/rs/zerolog/log"
)

type HandlerFunc func(net.Conn)

func (f HandlerFunc) HandleTCP(conn net.Conn) {
	f(conn)
}

type Handler interface {
	HandleTCP(net.Conn)
}

func TCP(h Handler) {
	addr := os.Getenv("ADDRESS")
	if addr == "" {
		addr = ":8080"
	}
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
