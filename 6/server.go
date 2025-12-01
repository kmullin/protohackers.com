package main

import (
	"io"
	"net"
	"time"

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

	Type         uint8          // to record what type of session this is, after identifying Camera or Dispatcher
	logger       zerolog.Logger // for session context aware logger
	heartbeating bool           // if we have a heartbeat running
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
			// It is an error for a client that has already identified itself as either a camera or a ticket dispatcher to send an IAmCamera message.
			switch ss.Type {
			case 0:
				ss.Type = message.MsgTypeIAmCamera
				ss.logger = ss.logger.With().Interface("camera", v).Logger()
			case message.MsgTypeIAmDispatcher:
				ss.Error("camera session tried to change to dispatcher")
				return
			case message.MsgTypeIAmCamera:
				ss.Error("camera session sent duplicate IAmCamera msg")
				return
			}

			ss.logger.Info().Msg("received message")
		case *message.IAmDispatcher:
			// It is an error for a client that has already identified itself as either a camera or a ticket dispatcher to send an IAmDispatcher message.
			switch ss.Type {
			case 0:
				ss.Type = message.MsgTypeIAmDispatcher
				ss.logger = ss.logger.With().Interface("dispatcher", v).Logger()
			case message.MsgTypeIAmCamera:
				ss.Error("dispatcher session tried to change to camera")
				return
			case message.MsgTypeIAmDispatcher:
				ss.Error("dispatcher session sent duplicate IAmDispatcher msg")
				return
			}

			ss.logger.Info().Msg("received message")
		case *message.Plate:
			// It is an error for a client that has not identified itself as a camera to send a Plate message.
			if ss.Type != message.MsgTypeIAmCamera {
				ss.Error("not a camera invalid plate msg")
				return
			}

			ss.logger.Info().Interface("plate", v).Msg("received message")
		case *message.WantHeartbeat:
			// It is an error for a client to send multiple WantHeartbeat messages on a single connection.
			if ss.heartbeating {
				ss.Error("heartbeat already requested")
				return
			}

			ss.logger.Info().Dur("want heartbeat", v.Interval).Msg("received message")
			ss.StartHeartbeat(v.Interval)
		default:
			ss.Error("unknown message")
			return
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

func (ss *Session) StartHeartbeat(d time.Duration) {
	ss.heartbeating = true
	ss.logger.Info().Dur("interval", d).Msg("starting heartbeat")
	// TODO: implement a heartbeat
}
