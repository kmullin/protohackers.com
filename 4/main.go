package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

const readTimeout = 100 * time.Millisecond

func main() {
	log.SetFlags(0)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	server := NewServer()
	server.Start(ctx)
	<-ctx.Done()
}
