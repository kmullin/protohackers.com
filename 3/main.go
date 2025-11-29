// Budget Chat
package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kmullin/protohackers.com/3/chat"
	"github.com/kmullin/protohackers.com/server"
)

func main() {
	var textLogging bool
	flag.BoolVar(&textLogging, "text", false, "turn on text logging")
	flag.Parse()

	if textLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	s := chat.NewServer(log.Logger)
	server.TCP(s)
}
