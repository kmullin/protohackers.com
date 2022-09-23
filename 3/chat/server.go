package chat

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"unicode"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	channels []string // placeholder
	users    []User

	logger zerolog.Logger // TODO: interface
}

func (s *Server) HandleTCP(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Printf("closed %s", conn.RemoteAddr())
	}()

	fmt.Fprintf(conn, "Welcome to budgetchat! What shall I call you?\n")

	scanner := bufio.NewScanner(conn)
	user := s.readUserName(scanner)
	log.Info().Interface("username", user).Msg("user joined")
}

// readUserName reads from the connection and returns the validated User
func (s *Server) readUserName(scanner *bufio.Scanner) User {
	scanner.Scan() // read a single line
	username := scanner.Text()

	if len(username) == 0 || !isAlphaNumeric(username) {
		log.Debug().Str("username", username).Msg("invalid username")
	}
	return User{username, time.Now()}
}

func isAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}
