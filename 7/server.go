package main

import (
	"context"
	"errors"
	"net"

	"github.com/kmullin/protohackers.com/7/message"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	ctx context.Context
	log zerolog.Logger
	sc  *SessionCache
}

func NewServer(ctx context.Context) *Server {
	return &Server{
		ctx: ctx,
		log: log.Logger,
		sc:  NewSessionCache(),
	}
}

func (s *Server) HandleUDP(conn net.PacketConn) {
	defer conn.Close()
	go func() {
		<-s.ctx.Done()
		conn.Close()
	}()

	for {
		buf := make([]byte, message.MaxSize)

		log := s.log
		n, addr, err := conn.ReadFrom(buf)
		if addr != nil {
			log = s.log.With().Stringer("addr", addr).Logger()
		}
		if err != nil {
			if errors.Is(err, net.ErrClosed) || s.ctx.Err() != nil {
				log.Err(s.ctx.Err()).Msg("reading from")
				return
			}
			log.Error().Err(err).Msg("reading from")
			continue
		}

		log.Debug().Int("len", n).Bytes("bytes", buf[:n]).Msg("read")

		// read our message, any invalid message returns an error
		msg, err := message.New(buf[:n])
		if err != nil {
			log.Error().Err(err).Msg("err creating msg")
			continue
		}

		switch m := msg.(type) {
		case *message.Connect:
			log.Info().Object("msg", m).Msg("recv connect msg")
			ss := s.sc.Get(m.SessionID)
			if ss == nil {
				ss = NewSession(m.SessionID, addr, conn, s.log)
				s.sc.Add(ss)
			}

			// will send to the same addr associated with session
			err := ss.Ack(0)
			if err != nil {
				log.Error().Err(err).Msg("writing ack")
				continue
			}
		case *message.Data:
			log.Info().Object("msg", m).Msg("recv data msg")

			ss := s.sc.Get(m.SessionID)
			// If the session is not open: send /close/SESSION/ and stop
			if ss == nil {
				log.Debug().Int("id", m.SessionID).
					Msg("couldnt find session, closing")
				err := ss.Close()
				if err != nil {
					log.Error().Err(err).Msg("err writing close msg")
				}
				continue
			}
			err := ss.InsertData(m)
			if err != nil {
				log.Error().Err(err).Msg("adding data to session")
			}
		case *message.Ack:
			log.Info().Object("msg", m).Msg("recv ack msg")

			ss := s.sc.Get(m.SessionID)
			// If the session is not open: send /close/SESSION/ and stop
			if ss == nil {
				log.Debug().Int("id", m.SessionID).
					Msg("couldnt find session, closing")
				err := ss.Close()
				if err != nil {
					log.Error().Err(err).Msg("err writing close msg")
				}
				continue
			}

			// If the LENGTH value is not larger than the largest LENGTH value in any ack message you've received on this session so far:
			//   do nothing and stop (assume it's a duplicate ack that got delayed).
			// If the LENGTH value is larger than the total amount of payload you've sent: the peer is misbehaving, close the session.
			// If the LENGTH value is smaller than the total amount of payload you've sent: retransmit all payload data after the first LENGTH bytes.
			// If the LENGTH value is equal to the total amount of payload you've sent: don't send any reply.

		case *message.Close:
			ss := s.sc.Get(m.SessionID)
			if ss != nil {
				err := ss.Close()
				if err != nil {
					log.Error().Err(err).Msg("err writing close msg")
				}
			}
		}
	}
}
