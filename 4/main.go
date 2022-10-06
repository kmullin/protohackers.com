package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var textLogging bool
	var debug bool
	flag.BoolVar(&textLogging, "text", false, "turn on text logging")
	flag.BoolVar(&debug, "d", false, "turn on debug logging")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if textLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	server := NewServer(log.Logger)
	server.Start(ctx)
	<-ctx.Done()
}
