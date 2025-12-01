package main

import (
	"io"
	"net"
	"os"
	"time"

	"github.com/kmullin/protohackers.com/6/message"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var msg io.WriterTo

	// client 1
	conn := MustDial()

	msg = &message.IAmCamera{Road: 123, Mile: 8, Limit: 60}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	msg = &message.WantHeartbeat{Interval: 2500 * time.Millisecond}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	msg = &message.Plate{Plate: "UN1X", Timestamp: time.Unix(0, 0).UTC()}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	time.Sleep(3 * time.Second)

	// client 2
	conn = MustDial()

	msg = &message.IAmCamera{Road: 123, Mile: 9, Limit: 60}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	msg = &message.Plate{Plate: "UN1X", Timestamp: time.Unix(45, 0).UTC()}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	// client 3
	conn = MustDial()

	msg = &message.IAmDispatcher{Roads: []uint16{123}}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

}

func MustDial() net.Conn {
	conn, err := net.Dial("tcp", os.Getenv("ADDRESS"))
	if err != nil {
		log.Fatal().Err(err).Msg("Dial failed")
	}

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	return conn
}
