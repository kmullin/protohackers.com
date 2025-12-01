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

	Type uint8

	logger zerolog.Logger
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
	// TODO: set a deadline for the connection and keep updating it after successful io
	ss := &Session{
		Conn:   conn,
		Type:   0,
		logger: s.logger.With().Stringer("remote", conn.RemoteAddr()).Logger(),
	}

	// tear down client connection after disconnect
	defer func() {
		if err := conn.Close(); err != nil {
			ss.logger.Err(err).Msg("disconnect")
		}
		ss.logger.Info().Msg("disconnected")
	}()

	ss.logger.Info().Msg("connected")
	for {
		msg, err := message.New(ss.Conn)
		if err != nil {
			if err == io.EOF {
				return
			}
			ss.logger.Err(err).Msg("parsing message")
		}

		switch v := msg.(type) {
		case *message.IAmCamera:
			switch ss.Type {
			case 0:
				ss.Type = message.MsgTypeIAmCamera
				ss.logger = ss.logger.With().Interface("camera", v).Logger()
			case message.MsgTypeIAmDispatcher:
				ss.Error("camera session tried to change to dispatcher")
				return
			}

			ss.logger.Info().Msg("received message")
		case *message.IAmDispatcher:
			switch ss.Type {
			case 0:
				ss.Type = message.MsgTypeIAmDispatcher
				ss.logger = ss.logger.With().Interface("dispatcher", v).Logger()
			case message.MsgTypeIAmCamera:
				ss.Error("dispatcher session tried to change to camera")
				return
			}

			ss.logger.Info().Msg("received message")
		case *message.Plate:
			ss.logger.Info().Interface("plate", v).Msg("received message")
		case *message.WantHeartbeat:
			ss.logger.Info().Dur("want heartbeat", v.Interval).Msg("received message")
		default:
			ss.logger.Error().Msg("unknown message")
		}
	}
}

// Error logs any errors and sends the client the same error message
func (ss *Session) Error(msg string) {
	ss.logger.Error().Msg(msg)

	e := &message.Error{Msg: msg}
	if _, err := e.WriteTo(ss); err != nil {
		ss.logger.Err(err)
	}
}
