package main

import (
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
	s.logger.Info().Stringer("client", conn.RemoteAddr()).Msg("connected")

	ss := new(Session)
	ss.Conn = conn

	msg, err := message.New(ss.Conn)
	if err != nil {
		s.logger.Err(err).Msg("parsing message")
	}
	s.logger.Debug().Interface("msg", msg)

	switch v := msg.(type) {
	case *message.IAmCamera:
		s.logger.Info().Interface("camera", v).Stringer("remote", conn.RemoteAddr()).Msg("connected")
	case *message.IAmDispatcher:
		s.logger.Info().Interface("dispatcher", v).Stringer("remote", conn.RemoteAddr()).Msg("connected")
	default:
		s.logger.Debug().Interface("type", v).Msg("unknown message")
	}
}
