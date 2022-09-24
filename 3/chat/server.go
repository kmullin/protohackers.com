package chat

import (
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	sessions []*Session
	mu       *sync.RWMutex

	logger zerolog.Logger // TODO: interface
}

func NewServer(logger zerolog.Logger) *Server {
	s := &Server{logger: logger, mu: new(sync.RWMutex)}
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
			s.logger.Info().Interface("users", s.sessions).Msg("currently connected")
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

	s.addSession(session)
	defer s.removeSession(session)
	s.logger.Info().Interface("session", session).Msg("user joined")

	err = session.ReadAll()
	if err != nil {
		s.logger.Err(err).Msg("")
		return
	}
}

func (s *Server) addSession(session *Session) {
	s.mu.Lock()
	s.sessions = append(s.sessions, session)
	s.mu.Unlock()
}

func (s *Server) removeSession(session *Session) {
	var i int
	var sesh *Session
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, sesh = range s.sessions {
		if sesh == session {
			break
		}
	}

	sessions := make([]*Session, 0)
	sessions = append(sessions, s.sessions[:i]...)
	sessions = append(sessions, s.sessions[i+1:]...)
	s.sessions = sessions
}
