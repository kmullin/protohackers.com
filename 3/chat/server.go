package chat

import (
	"net"
	"time"

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

	session, err := NewSession(conn)
	if err != nil {
		s.logger.Err(err).Msg("establishing session")
		return
	}
	s.logger.Info().Interface("session", session).Msg("user joined")
	time.Sleep(5 * time.Second)
}
