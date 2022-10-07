package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"git.kpmullin.com/kmullin/protocolhackers.com/4/database"
	reuse "github.com/libp2p/go-reuseport"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

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

	go s.logDbStatus(ctx)
	s.logger.Debug().Str("address", address).Msg("listening")
	return nil
}

func (s *server) handleUDP(ctx context.Context, conn net.PacketConn) {
	defer conn.Close()
	// TODO: something with context

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

		m := NewMessage(buf[:n])

		switch m.Type {
		case messageInsert:
			if m.Key == "version" {
				s.logger.Info().Msg("version request")
				continue
			}
			s.db.Insert(m.Key, m.Value)
			s.logger.Info().Str("type", "insert").Str("key", m.Key).Str("value", m.Value).Msg("")
		case messageRetrieve:
			var v string
			if m.Key == "version" {
				s.logger.Info().Msg("version request")
				v = version
			} else {
				v, _ = s.db.Retrieve(m.Key)
				s.logger.Info().Str("type", "retrieve").Str("key", m.Key).Str("value", v).Msg("")
			}
			response := fmt.Sprintf("%v=%v", m.Key, v)
			conn.WriteTo([]byte(response), addr)
		}
		s.logger.Debug().Int("bytes", n).Str("addr", addr.String()).Msg("done")
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
