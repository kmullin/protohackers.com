package chat

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

const welcomeMessage = "Welcome to budgetchat! What shall I call you?"

type Session struct {
	User       user
	TimeJoined time.Time
	recvC      chan message
	sendC      chan message

	conn    net.Conn
	scanner *bufio.Scanner
}

func NewSession(conn net.Conn) (*Session, error) {
	var err error
	var s Session
	s.TimeJoined = time.Now()
	s.recvC, s.sendC = make(chan message), make(chan message)
	s.conn = conn
	s.scanner = bufio.NewScanner(conn)
	s.scanner.Split(msgSplitFunc)

	s.User, err = s.getUserName()
	if err != nil {
		// inform client of bad username
		_, _ = s.WriteString("Invalid username (must be alphanumeric)")
		return nil, fmt.Errorf("reading username: %w", err)
	}
	return &s, nil
}

func (s *Session) getUserName() (user, error) {
	// start session with a welcome message
	_, err := s.WriteString(welcomeMessage)
	if err != nil {
		return user{}, err
	}
	msg, err := s.readMessage()
	if err != nil {
		return user{}, err
	}

	u := user{msg.String()}
	if !u.IsValid() {
		return user{}, fmt.Errorf("invalid username: %q", u.Name)
	}
	return u, nil
}

// readMessage reads a single message from the connection
func (s *Session) readMessage() (message, error) {
	if !s.scanner.Scan() {
		return message(""), fmt.Errorf("err scanning")
	}

	m := message(s.scanner.Text())
	if !m.IsValid() {
		return message(""), fmt.Errorf("invalid msg")
	}
	return m, nil
}

func (s *Session) ReadAll() error {
	for s.scanner.Scan() {
		m := message(s.scanner.Bytes())
		if !m.IsValid() {
			log.Debug().Interface("msg", m).Msg("invalid msg")
			continue
		}
		log.Info().Str("user", s.User.Name).Str("msg", m.String()).Msg("")
	}
	return s.scanner.Err()
}

// WriteString implements the io.StringWriter interface, and sends a newline terminated string to the active session
func (s *Session) WriteString(msg string) (int, error) {
	return fmt.Fprintln(s.conn, msg)
}

// msgSplitFunc is a bufio.SplitFunc that acts like ScanLines but ignores extra invalid lines on EOF
func msgSplitFunc(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF {
		// if we're already at EOF, we dont want any remaining data
		return 0, nil, nil
	}
	return bufio.ScanLines(data, atEOF)
}
