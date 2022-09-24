package chat

import (
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	users []User
	mu    *sync.RWMutex

	logger zerolog.Logger // TODO: interface
}

func NewServer(logger zerolog.Logger) *Server {
	s := &Server{logger: logger, users: make([]User, 0), mu: new(sync.RWMutex)}
	s.startStateLog()
	return s
}

func (s *Server) startStateLog() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			s.mu.RLock()
			s.logger.Info().Interface("users", s.users).Msg("currently connected")
			s.mu.RUnlock()
		}
	}()
}

func (s *Server) HandleTCP(conn net.Conn) {
	defer func() {
		conn.Close()
		s.logger.Printf("closed %s", conn.RemoteAddr())
	}()

	session, err := NewSession(conn)
	if err != nil {
		s.logger.Err(err).Msg("establishing session")
		return
	}

	s.addUser(session.User)
	defer s.removeUser(session.User)
	s.logger.Info().Interface("session", session).Msg("user joined")

	err = session.ReadAll()
	if err != nil {
		s.logger.Err(err).Msg("")
		return
	}
}

func (s *Server) addUser(user User) {
	s.mu.Lock()
	s.users = append(s.users, user)
	s.mu.Unlock()
}

func (s *Server) removeUser(user User) {
	var i int
	var u User
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, u = range s.users {
		if u == user {
			break
		}
	}

	users := make([]User, 0)
	users = append(users, s.users[:i]...)
	users = append(users, s.users[i+1:]...)
	s.users = users
}
