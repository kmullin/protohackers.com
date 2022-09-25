package chat

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	sessions []*Session
	mu       *sync.RWMutex

	msg chan message

	logger zerolog.Logger // TODO: interface
}

func NewServer(logger zerolog.Logger) *Server {
	s := Server{
		mu:     new(sync.RWMutex),
		msg:    make(chan message),
		logger: logger,
	}
	go s.startStateLog()
	go s.startRelay()
	return &s
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

	_ = s.announceSession(session)
	s.addSession(session)
	defer s.removeSession(session)
	s.logger.Info().Interface("session", session).Msg("user joined")

	err = session.ReadAll()
	if err != nil {
		s.logger.Err(err).Msg("")
		return
	}
}

func (s *Server) announceSession(session *Session) error {
	var users []string
	s.mu.RLock()
	for _, as := range s.sessions {
		_, err := as.WriteString(fmt.Sprintf("* %v has entered the room", session.User))
		if err != nil {
			s.logger.Err(err).Interface("session", s).Msg("writing to session")
		}
		users = append(users, as.User.Name)
	}
	s.mu.RUnlock()

	if len(users) > 0 {
		response := fmt.Sprintf("* Other peeps: %v", strings.Join(users, ", "))
		_, err := session.WriteString(response)
		if err != nil {
			s.logger.Err(err).Interface("session", session).Msg("writing to session")
		}
	}
	return nil
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

func (s *Server) startRelay() {

}

func (s *Server) startStateLog() {
	// FIXME: runs forever
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		s.mu.RLock()
		s.logger.Info().Interface("users", s.sessions).Msg("currently connected")
		s.mu.RUnlock()
	}
}
