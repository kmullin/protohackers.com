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

func (s *Server) HandleTCP(conn net.Conn) {
	// TODO: set a deadline for the connection and keep updating it after successful io
	ss := &Session{
		Conn:    conn,
		logger:  s.logger.With().Stringer("remote", conn.RemoteAddr()).Logger(),
		ticketC: make(chan message.Ticket, 32),
		hbDoneC: make(chan bool, 1),
	}

	// tear down client connection after disconnect
	defer func() {
		ss.hbDoneC <- true
		if err := conn.Close(); err != nil {
			ss.logger.Err(err).Msg("disconnect")
		}
		ss.logger.Info().Msg("disconnected")
	}()

	ss.logger.Info().Msg("connected")
	for {
		var msg message.Message

		select {
		case ticket := <-ss.ticketC:
			// dispatcher only
			_, err := ticket.WriteTo(conn)
			if err != nil {
				ss.logger.Error().Err(err).Msg("failed to write ticket")
				// XXX: return here?
			}
		default:
			// conn.SetDeadline(time.Now().Add(5 * time.Second))
			var err error
			msg, err = message.New(ss.Conn)
			if err != nil {
				if err != io.EOF {
					ss.logger.Err(err).Msg("parsing message")
				}
				return
			}
		}

		switch v := msg.(type) {
		case *message.IAmCamera:
			// It is an error for a client that has already identified itself as either a camera or a ticket dispatcher to send an IAmCamera message.
			switch {
			case ss.dispatcher != nil:
				ss.Error("dispatcher session tried to change to camera")
				return
			case ss.camera != nil:
				ss.Error("camera session sent duplicate IAmCamera msg")
				return
			case ss.camera == nil:
				ss.camera = v
				ss.logger = ss.logger.With().Interface("camera", v).Logger()
			}

			ss.logger.Info().Msg("new camera")
		case *message.IAmDispatcher:
			// It is an error for a client that has already identified itself as either a camera or a ticket dispatcher to send an IAmDispatcher message.
			switch {
			case ss.camera != nil:
				ss.Error("camera session tried to change to dispatcher")
				return
			case ss.dispatcher != nil:
				ss.Error("dispatcher session sent duplicate IAmDispatcher msg")
				return
			case ss.dispatcher == nil:
				ss.dispatcher = v
				ss.logger = ss.logger.With().Interface("dispatcher", v).Logger()
			}

			// XXX: need to start background handler for watching for tickets
			ss.logger.Info().Msg("new dispatcher")
			unsub := s.t.subs.Subscribe(v.Roads, ss.ticketC)
			defer unsub()
		case *message.Plate:
			// It is an error for a client that has not identified itself as a camera to send a Plate message.
			if ss.camera == nil {
				ss.Error("not a camera invalid plate msg")
				return
			}

			ss.logger.Info().Interface("plate", v).Msg("received plate from camera")
			s.t.Observe(v, ss.camera)
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
