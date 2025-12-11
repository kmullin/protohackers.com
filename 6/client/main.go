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

	/*
		t, err := time.Parse(time.RFC3339, "1970-01-09T14:21:58Z")
		if err != nil {
			log.Error().Err(err).Msg("failed to parse time")
		}

		fmt.Print(t.Unix())
		os.Exit(0)
	*/
	// client 1
	conn := MustDial()

	msg = &message.IAmCamera{Road: 2220, Mile: 419, Limit: 100}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	t, _ := time.Parse(time.RFC3339, "1970-01-02T04:50:40Z")
	msg = &message.Plate{Plate: "UP19RHG", Timestamp: t}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	// 	msg = &message.IAmDispatcher{Roads: []uint16{123}}
	// 	if _, err := msg.WriteTo(conn); err != nil {
	// 		log.Error().Err(err).Msg("write failed")
	// 	}

	// client 2
	conn = MustDial()

	msg = &message.IAmCamera{Road: 2220, Mile: 429, Limit: 100}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	msg = &message.Plate{Plate: "UP19RHG", Timestamp: tParse("1970-01-02T04:55:40Z")}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	// client 3
	conn = MustDial()

	msg = &message.IAmDispatcher{Roads: []message.RoadID{2220}}
	if _, err := msg.WriteTo(conn); err != nil {
		log.Error().Err(err).Msg("write failed")
	}

	time.Sleep(3 * time.Second)
}

func MustDial() net.Conn {
	conn, err := net.Dial("tcp", os.Getenv("ADDRESS"))
	if err != nil {
		log.Fatal().Err(err).Msg("Dial failed")
	}

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	return conn
}

func tParse(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
