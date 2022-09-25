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

	msgs chan message

	logger zerolog.Logger // TODO: interface
}

func NewServer(logger zerolog.Logger) *Server {
	s := Server{
		mu:     new(sync.RWMutex),
		msgs:   make(chan message),
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

	s.announceJoin(session)
	s.addSession(session)
	defer func() {
		s.removeSession(session)
		s.announceLeave(session)
	}()
	s.logger.Info().Interface("session", session).Msg("user joined")

	err = session.ReadAll(s.msgs)
	if err != nil {
		s.logger.Err(err).Msg("")
		return
	}
}

// announceJoin announces session to all current active sessions
func (s *Server) announceJoin(session *Session) {
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

	response := fmt.Sprintf("* Other peeps: %v", strings.Join(users, ", "))
	_, err := session.WriteString(response)
	if err != nil {
		s.logger.Err(err).Interface("session", session).Msg("writing to session")
	}
}

// announceLeave announces to all other sessions the leaving of another user
func (s *Server) announceLeave(session *Session) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, as := range s.sessions {
		_, err := as.WriteString(fmt.Sprintf("* %v has left the room", session.User))
		if err != nil {
			s.logger.Err(err).Interface("session", s).Msg("writing to session")
		}
	}
}

// addSession adds a session to the global session state
func (s *Server) addSession(session *Session) {
	s.mu.Lock()
	s.sessions = append(s.sessions, session)
	s.mu.Unlock()
}

// removeSession removes the session from the global session state
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

// startRelay starts a loop on reading the channel given to each sesssion
func (s *Server) startRelay() {
	for msg := range s.msgs {
		s.mu.RLock()
		for _, sesh := range s.sessions {
			if sesh != msg.session {
				m := fmt.Sprintf("[%s] %s", msg.session.User.Name, msg.msg)
				_, _ = sesh.WriteString(m)
			}
		}
		s.mu.RUnlock()
	}
}

// startStateLog starts an infinite loop that prints the current sesssions every 5 seconds
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
