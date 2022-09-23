package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"git.kpmullin.com/kmullin/protocolhackers.com/3/chat"
	"git.kpmullin.com/kmullin/protocolhackers.com/server"
)

func main() {
	var jsonLogging bool
	flag.BoolVar(&jsonLogging, "json", false, "turn on json logging")
	flag.Parse()

	if !jsonLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	s := chat.NewServer(log.Logger)
	server.TCP(s)
}
