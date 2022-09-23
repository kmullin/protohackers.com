package chat

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"unicode"

	"github.com/rs/zerolog"
)

type Server struct {
	channels []string // placeholder
	users    []User

	logger zerolog.Logger // TODO: interface
}

func NewServer(logger zerolog.Logger) *Server {
	return &Server{logger: logger}
}

func (s *Server) HandleTCP(conn net.Conn) {
	defer func() {
		conn.Close()
		s.logger.Printf("closed %s", conn.RemoteAddr())
	}()

	fmt.Fprintf(conn, "Welcome to budgetchat! What shall I call you?\n")

	scanner := bufio.NewScanner(conn)
	user := s.readUserName(scanner)
	s.logger.Info().Interface("username", user).Msg("user joined")
}

// readUserName reads from the connection and returns the validated User
func (s *Server) readUserName(scanner *bufio.Scanner) User {
	scanner.Scan() // read a single line
	username := scanner.Text()

	if len(username) == 0 || !isAlphaNumeric(username) {
		s.logger.Debug().Str("username", username).Msg("invalid username")
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
