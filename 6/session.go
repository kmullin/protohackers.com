package main

import (
	"net"
	"time"

	"github.com/kmullin/protohackers.com/6/message"
	"github.com/rs/zerolog"
)

type Session struct {
	net.Conn

	log zerolog.Logger // for session context aware logger

	camera     *message.IAmCamera     // to record camera information
	dispatcher *message.IAmDispatcher // to record dispatcher information
	ticketC    <-chan message.Ticket

	heartbeating bool // if we have a heartbeat running

	done chan struct{} // used to signal disconnect and to stop any heartbeating
}

func NewSession(conn net.Conn, logger zerolog.Logger) *Session {
	ss := &Session{
		Conn: conn,
		log:  logger.With().Stringer("remote", conn.RemoteAddr()).Logger(),
		done: make(chan struct{}),
	}
	ss.log.Info().Msg("connected")
	return ss
}

// Error logs any errors and sends the client the same error message
func (ss *Session) Error(msg string) {
	ss.log.Error().Msg(msg)

	e := &message.Error{Msg: msg}
	if _, err := e.WriteTo(ss); err != nil {
		ss.log.Error().Err(err).Msg("writing error")
	}
}

// StartHeartbeat starts sending a heartbeat message to the client at every interval d
func (ss *Session) StartHeartbeat(d time.Duration) {
	ss.heartbeating = true
	if d == 0 {
		return
	}
	ss.log.Info().Dur("interval", d).Msg("starting heartbeat")

	ticker := time.NewTicker(d)
	go func() {
		defer ticker.Stop()

		hb := &message.Heartbeat{}
		for {
			select {
			case <-ss.done:
				return
			case <-ticker.C:
				if _, err := hb.WriteTo(ss); err != nil {
					ss.log.Err(err).Msg("writing heartbeat")
				}
			}
		}
	}()
}

func (ss *Session) StartTicketing() {
	go func() {
		for {
			select {
			case ticket := <-ss.ticketC:
				_, err := ticket.WriteTo(ss)
				if err != nil {
					ss.log.Error().Err(err).Msg("failed to write ticket")
					return
				}
				ss.log.Info().Interface("ticket", ticket).Msg("sent ticket")
			case <-ss.done:
				ss.log.Info().Msg("stop ticketing")
				return
			}
		}
	}()
}

func (ss *Session) IsFresh() bool {
	return !ss.IsCamera() && !ss.IsDispatcher()
}

func (ss *Session) IsDispatcher() bool {
	return ss.dispatcher != nil
}

func (ss *Session) IsCamera() bool {
	return ss.camera != nil
}

func (ss *Session) ReadMsg() (message.Message, error) {
	return message.New(ss.Conn)
}

func (ss *Session) Close() {
	// tear down client connection after disconnect
	close(ss.done)
	if err := ss.Conn.Close(); err != nil {
		ss.log.Error().Err(err).Msg("disconnecting")
	}
	ss.log.Info().Msg("disconnected")
}
