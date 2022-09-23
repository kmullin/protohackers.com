package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"git.kpmullin.com/kmullin/protocolhackers.com/3/chat"
	"git.kpmullin.com/kmullin/protocolhackers.com/server"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	handler := &chat.Server{}
	server.TCP(handler)
}
