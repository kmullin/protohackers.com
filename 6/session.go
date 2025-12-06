package main

import (
	"net"
	"time"

	"github.com/kmullin/protohackers.com/6/message"
	"github.com/rs/zerolog"
)

type Session struct {
	net.Conn

	logger zerolog.Logger // for session context aware logger

	camera     *message.IAmCamera     // to record camera information
	dispatcher *message.IAmDispatcher // to record dispatcher information
	ticketC    chan message.Ticket

	heartbeating bool      // if we have a heartbeat running
	hbDoneC      chan bool // used to signal disconnect and to stop any heartbeating
}

// Error logs any errors and sends the client the same error message
func (ss *Session) Error(msg string) {
	ss.logger.Error().Msg(msg)

	e := &message.Error{Msg: msg}
	if _, err := e.WriteTo(ss); err != nil {
		ss.logger.Err(err).Msg("writing error")
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
			case <-ss.hbDoneC:
				ss.logger.Info().Msg("stopping heartbeat")
				return
			case <-ticker.C:
				ss.logger.Info().Msg("sending heartbeat")

				if _, err := hb.WriteTo(ss); err != nil {
					ss.logger.Err(err).Msg("writing heartbeat")
				}
			}
		}
	}()
}
