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
		s.logger.Printf("Listening on %v...", port)
		conn, err := reuse.ListenPacket("udp", port)
		if err != nil {
			s.logger.Printf("unable to listen: %v", err)
		}
		go s.handleUDP(ctx, conn)
	}
	go s.logDbStatus(ctx)
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
				s.logger.Printf("UDP read err: %v", err)
			}
			continue
		}

		if n == bufSize {
			s.logger.Printf("invalid message size %v from %v, ignoring", n, addr)
			drained, err := drain(conn)
			if err != nil {
				s.logger.Printf("err draining: %v", err)
			}
			s.logger.Printf("drained %v bytes from %v", drained, addr)
			continue
		}

		s.logger.Printf("%v bytes received from %v: % x", n, addr, buf[:n])

		m := NewMessage(buf[:n])
		switch m.Type {
		case messageInsert:
			s.db.Insert(m.Key, m.Value)
		case messageRetrieve:
		}
		s.logger.Printf("msg %+v", m)
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
			f := s.db.Status()
			s.logger.Debug().Interface("database", f).Msg("")
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
