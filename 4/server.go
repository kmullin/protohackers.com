package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/kmullin/protohackers.com/4/database"
	reuse "github.com/libp2p/go-reuseport"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

// All requests and responses must be shorter than 1000 bytes.
const bufSize = 1000

const version = "kmullin's terrible K/V Store 420"

const readTimeout = 100 * time.Millisecond

type server struct {
	db     *database.Db
	logger zerolog.Logger
}

func NewServer(log zerolog.Logger) *server {
	return &server{database.NewDB(), log}
}

func (s *server) Start(ctx context.Context, address string) error {
	var g errgroup.Group
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		g.Go(func() error {
			conn, err := reuse.ListenPacket("udp", address)
			if err != nil {
				return err
			}
			go s.handleUDP(ctx, conn)
			return nil
		})
	}

	// Wait for all shared listeners to come up
	if err := g.Wait(); err != nil {
		return err
	}

	s.logger.Debug().Str("address", address).Msg("listening")
	return nil
}

func (s *server) handleUDP(ctx context.Context, conn net.PacketConn) {
	defer conn.Close()
	// TODO: something with context

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
			if err := drain(conn); err != nil {
				s.logger.Error().Err(err).Msg("err draining")
			}
			s.logger.Debug().Int("bytes", n).Str("addr", addr.String()).Msg("invalid message")
			continue
		}

		m := NewMessage(buf[:n])
		switch {
		case m.Key == "version":
			if m.Type == messageRetrieve {
				s.logger.Info().Msg("version request")
				conn.WriteTo(responseMsg(m.Key, version), addr)
			}
		case m.Type == messageInsert:
			s.logger.Info().Str("type", "insert").Str("key", m.Key).Str("value", m.Value).Send()
			s.db.Insert(m.Key, m.Value)
			s.logger.Debug().Int("entries", s.db.Entries()).Msg("new entry")
		case m.Type == messageRetrieve:
			v, _ := s.db.Retrieve(m.Key)
			conn.WriteTo(responseMsg(m.Key, v), addr)
			s.logger.Info().Str("type", "retrieve").Str("key", m.Key).Str("value", v).Send()
		}
		s.logger.Debug().Int("bytes", n).Str("addr", addr.String()).Msg("done")
	}
}

func readFrom(conn net.PacketConn, buf []byte) (int, net.Addr, error) {
	conn.SetDeadline(time.Now().Add(readTimeout))
	return conn.ReadFrom(buf)
}

// drain reads any remaining bytes in the PacketConn, returns total bytes read, and any error encountered
func drain(conn net.PacketConn) error {
	var n int
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
	}
	return err
}

func responseMsg(k, v string) []byte {
	return []byte(fmt.Sprintf("%v=%v", k, v))
}
