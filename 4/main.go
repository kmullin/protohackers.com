// Unusual Database Program
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
	flag.BoolVar(&textLogging, "text", false, "turn on text logging")
	flag.Parse()

	if textLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	addr := ":8080"
	_, isFlyIO := os.LookupEnv("FLY_ALLOC_ID")
	if isFlyIO {
		// https://fly.io/docs/app-guides/udp-and-tcp/
		addr = "fly-global-services:65534"
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	server := NewServer(log.Logger)
	err := server.Start(ctx, addr)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to listen")
	}
	<-ctx.Done()
}
