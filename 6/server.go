package main

import (
	"net"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	//observations []observation
	//mu     *sync.RWMutex
	logger zerolog.Logger
}

type Camera struct {
	Road       string
	Location   string
	SpeedLimit string
}

type Picture struct {
	Plate     string
	Timestamp time.Time
}

func (s *Server) HandleTCP(conn net.Conn) {
	// tear down client connection after disconnect
	defer func() {
		conn.Close()
		s.logger.Info().Stringer("client", conn.RemoteAddr()).Msg("disconnected")
	}()
	s.logger.Info().Stringer("client", conn.RemoteAddr()).Msg("connected")

}
