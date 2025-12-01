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

	logger zerolog.Logger // for session context aware logger
	Type   uint8          // to record what type of session this is, after identifying Camera or Dispatcher

	heartbeating bool      // if we have a heartbeat running
	doneC        chan bool // used to signal disconnect and to stop any heartbeating
}

func (s *Server) HandleTCP(conn net.Conn) {
	// TODO: set a deadline for the connection and keep updating it after successful io
	ss := &Session{
		Conn:   conn,
		Type:   0,
		logger: s.logger.With().Stringer("remote", conn.RemoteAddr()).Logger(),
		doneC:  make(chan bool, 1),
	}

	// tear down client connection after disconnect
	defer func() {
		ss.doneC <- true
		if err := conn.Close(); err != nil {
			ss.logger.Err(err).Msg("disconnect")
		}
		ss.logger.Info().Msg("disconnected")
	}()

	ss.logger.Info().Msg("connected")
	for {
		// conn.SetDeadline(time.Now().Add(5 * time.Second))
		msg, err := message.New(ss.Conn)
		if err != nil {
			if err != io.EOF {
				ss.logger.Err(err).Msg("parsing message")
			}
			return
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

			ss.logger.Info().Msg("new camera")
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

			ss.logger.Info().Msg("new dispatcher")
		case *message.Plate:
			// It is an error for a client that has not identified itself as a camera to send a Plate message.
			if ss.Type != message.MsgTypeIAmCamera {
				ss.Error("not a camera invalid plate msg")
				return
			}

			ss.logger.Info().Interface("plate", v).Msg("received plate from camera")
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

// StartHeartbeat starts sending a heartbeat message to the client at every interval d
func (ss *Session) StartHeartbeat(d time.Duration) {
	ss.heartbeating = true
	if d == 0 {
		return
	}
	ss.logger.Info().Dur("interval", d).Msg("starting heartbeat")

	ticker := time.NewTicker(d)
	go func() {
		defer ticker.Stop()

		hb := &message.Heartbeat{}
		for {
			select {
			case <-ss.doneC:
				ss.logger.Info().Msg("stopping heartbeat")
				return
			case <-ticker.C:
				ss.logger.Info().Msg("sending heartbeat")

				if _, err := hb.WriteTo(ss); err != nil {
					ss.logger.Err(err)
				}
			}
		}
	}()
}
