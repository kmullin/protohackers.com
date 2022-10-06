package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const readTimeout = 100 * time.Millisecond

func main() {
	var textLogging bool
	flag.BoolVar(&textLogging, "text", false, "turn on text logging")
	flag.Parse()

	if textLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	server := NewServer(log.Logger)
	server.Start(ctx)
	<-ctx.Done()
}
