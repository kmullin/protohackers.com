package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kmullin/protohackers.com/5/proxy"
	"github.com/kmullin/protohackers.com/server"
)

func main() {
	var textLogging bool
	var upstream string
	flag.BoolVar(&textLogging, "text", false, "turn on text logging")
	flag.StringVar(&upstream, "upstream", "", "use different upstream")
	flag.Parse()

	if textLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	p := proxy.NewServer(log.Logger)
	if upstream != "" {
		p.Upstream = upstream
	}
	server.TCP(p)
}
