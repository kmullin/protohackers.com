package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"time"

	reuse "github.com/libp2p/go-reuseport"
)

func main() {
	log.SetFlags(0)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	port := ":8080"

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		log.Printf("Listening on %v...", port)
		conn, err := reuse.ListenPacket("udp", port)
		if err != nil {
			log.Fatalf("unable to listen: %v", err)
		}
		go handleUDP(ctx, conn)
	}

	<-ctx.Done()
}

func handleUDP(ctx context.Context, conn net.PacketConn) {
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
			drain(conn)
			continue
		}

		log.Printf("%v bytes received from %v: % x", n, addr, buf[:n])
	}
}

func readFrom(conn net.PacketConn, buf []byte) (int, net.Addr, error) {
	setTimeout(conn)
	return conn.ReadFrom(buf[0:])
}

func drain(conn net.PacketConn) {
	var n, total int
	var err error

	n = 1
	buf := make([]byte, 4096)
	for n > 0 {
		n, _, err = readFrom(conn, buf[0:])
		if err != nil {
			if !errors.Is(err, os.ErrDeadlineExceeded) {
				log.Print(err)
			}
			break
		}
		total += n
	}
	log.Printf("drained %v bytes", total)
}

func setTimeout(conn net.PacketConn) {
	conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
}
