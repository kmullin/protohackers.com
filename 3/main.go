package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"unicode"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"git.kpmullin.com/kmullin/protocolhackers.com/server"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	server.TCP(server.HandlerFunc(handleConn))
}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Printf("closed %s", conn.RemoteAddr())
	}()
	fmt.Fprintf(conn, "Welcome to budgetchat! What shall I call you?\n")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 || !IsAlphaNumeric(t) {
			log.Debug().Str("username", t).Msg("invalid username")
			break
		}
		log.Info().Str("username", t).Msg("user joined")
	}

	if err := scanner.Err(); err != nil {
		log.Printf("reading standard input: %v", err)
	}
}

func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}
