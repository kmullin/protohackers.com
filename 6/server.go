package main

import (
	"io"
	"net"

	"github.com/kmullin/protohackers.com/6/message"
	"github.com/rs/zerolog"
)

type Server struct {
	//observations []observation
	//mu     *sync.RWMutex
	logger zerolog.Logger
}

type Session struct {
	net.Conn

	Type int // the session Type
}

//	type Camera struct {
//		Road       string
//		Location   string
//		SpeedLimit string
//	}
//
//	type Picture struct {
//		Plate     string
//		Timestamp time.Time
//	}

func (s *Server) HandleTCP(conn net.Conn) {
	// tear down client connection after disconnect
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Err(err).Stringer("client", conn.RemoteAddr()).Msg("disconnect")
		}
		s.logger.Info().Stringer("client", conn.RemoteAddr()).Msg("disconnected")
	}()

	// TODO: set a deadline for the connection and keep updating it after successful io
	ss := new(Session)
	ss.Conn = conn

	s.logger.Info().Stringer("remote", conn.RemoteAddr()).Msg("connected")
	for {
		msg, err := message.New(ss.Conn)
		if err != nil {
			if err == io.EOF {
				return
			}
			s.logger.Err(err).Msg("parsing message")
		}

		switch v := msg.(type) {
		case *message.IAmCamera:
			s.logger.Info().Interface("camera", v).Stringer("remote", conn.RemoteAddr()).Msg("received message")
		case *message.IAmDispatcher:
			s.logger.Info().Interface("dispatcher", v).Stringer("remote", conn.RemoteAddr()).Msg("received message")
		}
	}
}
