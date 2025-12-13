package main

import (
	"flag"
	"os"

	"github.com/kmullin/protohackers.com/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var textLogging bool
	flag.BoolVar(&textLogging, "text", false, "turn on text logging")
	flag.Parse()

	if textLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	s := NewServer(log.Logger)
	server.UDP(s)
}
