package main

import (
	"io"
	"net"

	"github.com/kmullin/protohackers.com/6/message"
	"github.com/rs/zerolog"
)

type Server struct {
	logger zerolog.Logger
	t      *Ticketer
}

func NewServer(logger zerolog.Logger) *Server {
	return &Server{
		logger: logger,
		t:      NewTicketer(logger),
	}
}

func (s *Server) HandleTCP(conn net.Conn) {
	ss := NewSession(conn, s.logger)
	defer ss.Close()

	for {
		msg, err := ss.ReadMsg()
		if err != nil {
			if err != io.EOF {
				ss.Error("illegal message")
			}
			return
		}

		switch v := msg.(type) {
		case *message.IAmCamera:
			// It is an error for a client that has already identified itself as either a camera or a ticket dispatcher to send an IAmCamera message.
			switch {
			case ss.IsDispatcher():
				ss.Error("dispatcher session tried to change to camera")
				return
			case ss.IsCamera():
				ss.Error("camera session sent duplicate IAmCamera msg")
				return
			case ss.IsFresh():
				ss.camera = v
				ss.log = ss.log.With().Interface("camera", v).Logger()
			}

			ss.log.Info().Msg("new camera")

			// create our channel for tickets on this road
			_ = s.t.subs.NewRoad(v.Road)
		case *message.IAmDispatcher:
			// It is an error for a client that has already identified itself as either a camera or a ticket dispatcher to send an IAmDispatcher message.
			switch {
			case ss.IsCamera():
				ss.Error("camera session tried to change to dispatcher")
				return
			case ss.IsDispatcher():
				ss.Error("dispatcher session sent duplicate IAmDispatcher msg")
				return
			case ss.IsFresh():
				ss.dispatcher = v
				ss.log = ss.log.With().Interface("dispatcher", v).Logger()
			}

			// XXX: need to start background handler for watching for tickets
			ss.log.Info().Msg("new dispatcher")
			ss.ticketC = s.t.subs.Subscribe(ss.done, v.Roads)
			ss.StartTicketing()
		case *message.Plate:
			// It is an error for a client that has not identified itself as a camera to send a Plate message.
			if !ss.IsCamera() {
				ss.Error("not a camera invalid plate msg")
				return
			}

			s.t.Observe(v, ss.camera)
		case *message.WantHeartbeat:
			// It is an error for a client to send multiple WantHeartbeat messages on a single connection.
			if ss.heartbeating {
				ss.Error("heartbeat already requested")
				return
			}

			ss.log.Info().Dur("want heartbeat", v.Interval).Msg("received message")
			ss.StartHeartbeat(v.Interval)
		case *message.Error:
			ss.Error("server recevied a server error message")
			return
		case *message.Ticket:
			ss.Error("server recevied a server ticket message")
			return
		case *message.Heartbeat:
			ss.Error("server recevied a server heartbeat message")
			return
		default:
			ss.Error("unknown message")
			return
		}
	}
}
