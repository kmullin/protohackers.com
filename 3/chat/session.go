package chat

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Session struct {
	id         uint // ?
	User       User
	TimeJoined time.Time
	scanner    *bufio.Scanner
}

func NewSession(conn net.Conn) (*Session, error) {
	var err error
	var s Session
	s.TimeJoined = time.Now()

	_, err = fmt.Fprintln(conn, "Welcome to budgetchat! What shall I call you?")
	if err != nil {
		return nil, err
	}

	// setup session scanner for line reads
	s.scanner = bufio.NewScanner(conn)
	s.User, err = s.readUserName()
	if err != nil {
		return nil, fmt.Errorf("reading username: %w", err)
	}
	return &s, nil
}

func (s *Session) readUserName() (User, error) {
	s.scanner.Scan()
	user := User{s.scanner.Text()}
	if !user.IsValid() {
		return User{}, fmt.Errorf("invalid username: %q", user.Name)
	}
	return user, nil
}

func (s *Session) ReadAll() error {
	for s.scanner.Scan() {

	}
	if err := s.scanner.Err(); err != nil {
		return fmt.Errorf("reading client: %w", err)
	}
	return nil
}
