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

	scanner *bufio.Scanner
	conn    net.Conn
}

func NewSession(conn net.Conn) (*Session, error) {
	var err error
	var s Session
	s.TimeJoined = time.Now()
	s.recvC, s.sendC = make(chan Message), make(chan Message)
	s.scanner = bufio.NewScanner(conn)
	s.scanner.Split(splitFunc)
	s.conn = conn

	err = s.writeMsg("Welcome to budgetchat! What shall I call you?")
	if err != nil {
		return nil, err
	}

	s.User, err = s.readUserName()
	if err != nil {
		// inform client of bad username
		_ = s.writeMsg("Invalid username (must be alphanumeric)")
		return nil, fmt.Errorf("reading username: %w", err)
	}
	return &s, nil
}

func (s *Session) readUserName() (User, error) {
	if !s.scanner.Scan() {
		return User{}, fmt.Errorf("couldnt scan for username")
	}

	user := User{s.scanner.Text()}
	if !user.IsValid() {
		return User{}, fmt.Errorf("invalid username: %q", user.Name)
	}
	return user, nil
}

func (s *Session) ReadAll() error {
	for s.scanner.Scan() {
		log.Info().Str("user", s.User.Name).Str("msg", s.scanner.Text()).Msg("")
	}
	if err := s.scanner.Err(); err != nil {
		log.Err(err).Msg("")
	}
	return nil
}

func (s *Session) writeMsg(msg string) error {
	_, err := fmt.Fprintln(s.conn, msg)
	return err
}
