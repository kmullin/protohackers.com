package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	reuse "github.com/libp2p/go-reuseport"
)

type server struct {
	db *db
}

func NewServer() *server {
	return &server{NewDB()}
}

func (s *server) Start(ctx context.Context) error {
	port := ":8080"
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		log.Printf("Listening on %v...", port)
		conn, err := reuse.ListenPacket("udp", port)
		if err != nil {
			log.Fatalf("unable to listen: %v", err)
		}
		go s.handleUDP(ctx, conn)
	}
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
				log.Printf("UDP read err: %v", err)
			}
			continue
		}

		if n == bufSize {
			log.Printf("invalid message size %v from %v, ignoring", n, addr)
			drained, err := drain(conn)
			if err != nil {
				log.Printf("err draining: %v", err)
			}
			log.Printf("drained %v bytes from %v", drained, addr)
			continue
		}

		log.Printf("%v bytes received from %v: % x", n, addr, buf[:n])
		m := msg(buf[:n])
		k, v := m.KV()
		log.Printf("msg %q = %q", k, v)
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
