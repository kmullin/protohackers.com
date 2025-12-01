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

	conn, err := net.Dial("tcp", os.Getenv("ADDRESS"))
	if err != nil {
		log.Fatal().Err(err).Msg("Dial failed")
	}

	conn.SetDeadline(time.Now().Add(3 * time.Second))

	var msg io.WriterTo

	msg = &message.IAmCamera{Road: 123, Mile: 456, Limit: 100}
	_, err = msg.WriteTo(conn)
	if err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	// msg = &message.IAmDispatcher{Roads: []uint16{66}}
	// _, err = msg.WriteTo(conn)
	// if err != nil {
	// 	log.Error().Err(err).Msg("write failed")
	// }

	msg = &message.Plate{Plate: "thing", Timestamp: time.Now().UTC()}
	_, err = msg.WriteTo(conn)
	if err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	msg = &message.WantHeartbeat{1 * time.Second}
	_, err = msg.WriteTo(conn)
	if err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	n, err := io.Copy(os.Stderr, conn)
	if err != nil {
		log.Error().Err(err).Msg("read failed")
	}

	log.Debug().Int64("inside", n)
}
