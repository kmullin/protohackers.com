package log

import (
	"flag"
	"log"
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() {
	var textLogging bool
	flag.BoolVar(&textLogging, "text", false, "turn on text logging")
	flag.Parse()

	if textLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
