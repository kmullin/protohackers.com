package main

import (
	"context"
	"errors"
	"net"
	"os"
	"runtime"
	"time"

	"git.kpmullin.com/kmullin/protocolhackers.com/4/database"
	reuse "github.com/libp2p/go-reuseport"
	"github.com/rs/zerolog"
)

const readTimeout = 100 * time.Millisecond

type server struct {
	db     *database.Db
	logger zerolog.Logger
}

func NewServer(log zerolog.Logger) *server {
	return &server{database.NewDB(), log}
}

func (s *server) Start(ctx context.Context) error {
	port := ":8080"
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		conn, err := reuse.ListenPacket("udp", port)
		if err != nil {
			s.logger.Error().Err(err).Msg("unable to listen")
		}
		go s.handleUDP(ctx, conn)
	}
	go s.logDbStatus(ctx)
	s.logger.Debug().Str("port", port).Msg("listening")
	return nil
}

func (s *server) handleUDP(ctx context.Context, conn net.PacketConn) {
	defer conn.Close()

	// All requests and responses must be shorter than 1000 bytes.
	const bufSize = 1000
	buf := make([]byte, bufSize)
	for {
		n, addr, err := readFrom(conn, buf[0:])
		if err != nil {
			if !errors.Is(err, os.ErrDeadlineExceeded) {
				s.logger.Error().Err(err).Msg("UDP read err")
			}
			continue
		}

		if n == bufSize {
			drained, err := drain(conn)
			if err != nil {
				s.logger.Error().Err(err).Msg("err draining")
			}
			s.logger.Debug().Int("bytes", n+drained).Str("addr", addr.String()).Msg("invalid message")
			continue
		}

		s.logger.Debug().Int("bytes", n).Str("addr", addr.String()).Msg("")

		m := NewMessage(buf[:n])
		switch m.Type {
		case messageInsert:
			s.db.Insert(m.Key, m.Value)
		case messageRetrieve:
			v, _ := s.db.Retrieve(m.Key)
			s.logger.Debug().Str("value", v).Msg("retrieved")
		}
		s.logger.Info().Interface("type", m.Type).Str("key", m.Key).Str("value", m.Value).Msg("")
	}
}

func (s *server) logDbStatus(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d := s.db.Status()
			s.logger.Debug().Interface("database", d).Int("size", len(d)).Msg("")
		}
	}
}

func readFrom(conn net.PacketConn, buf []byte) (int, net.Addr, error) {
	conn.SetDeadline(time.Now().Add(readTimeout))
	return conn.ReadFrom(buf[0:])
}

// drain reads any remaining bytes in the PacketConn, returns total bytes read, and any error encountered
func drain(conn net.PacketConn) (int, error) {
	var n, total int
	var err error

	n = 1
	buf := make([]byte, 4096)
	for n > 0 {
		n, _, err = readFrom(conn, buf[0:])
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				err = nil
			}
			break
		}
		total += n
	}
	return total, err
}
