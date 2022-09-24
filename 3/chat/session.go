package chat

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

type Session struct {
	User       User
	TimeJoined time.Time
	recvC      chan Message
	sendC      chan Message

	conn    net.Conn
	scanner *bufio.Scanner
}

func NewSession(conn net.Conn) (*Session, error) {
	var err error
	var s Session
	s.TimeJoined = time.Now()
	s.recvC, s.sendC = make(chan Message), make(chan Message)
	s.conn = conn
	s.scanner = bufio.NewScanner(conn)
	s.scanner.Split(splitFunc)

	s.User, err = s.readUserName()
	if err != nil {
		// inform client of bad username
		_, _ = s.WriteString("Invalid username (must be alphanumeric)")
		return nil, fmt.Errorf("reading username: %w", err)
	}
	return &s, nil
}

func (s *Session) readUserName() (User, error) {
	// start session with a welcome message
	_, err := s.WriteString("Welcome to budgetchat! What shall I call you?")
	if err != nil {
		return User{}, err
	}
	msg, err := s.readMessage()
	if err != nil {
		return User{}, err
	}

	user := User{msg.String()}
	if !user.IsValid() {
		return User{}, fmt.Errorf("invalid username: %q", user.Name)
	}
	return user, nil
}

// readMessage reads a single message from the connection
func (s *Session) readMessage() (Message, error) {
	if !s.scanner.Scan() {
		return Message(""), fmt.Errorf("err scanning")
	}

	m := Message(s.scanner.Text())
	if !m.isValid() {
		return Message(""), fmt.Errorf("invalid msg")
	}
	return m, nil
}

func (s *Session) ReadAll() error {
	for {
		msg, err := s.readMessage()
		if err != nil {
			log.Err(err).Msg("")

		}
		log.Info().Str("user", s.User.Name).Str("msg", msg.String()).Msg("")
	}
	return nil
}

// WriteString implements the io.StringWriter interface, and sends a newline terminated string to the active session
func (s *Session) WriteString(msg string) (int, error) {
	return fmt.Fprintln(s.conn, msg)
}
