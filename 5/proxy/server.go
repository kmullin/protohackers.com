package proxy

import (
	"bufio"
	"fmt"
	"net"

	"github.com/rs/zerolog"
)

const upstream = "chat.protohackers.com:16963"

type Server struct {
	logger   zerolog.Logger
	Upstream string
}

func NewServer(logger zerolog.Logger) *Server {
	return &Server{logger: logger, Upstream: upstream}
}

// HandleTCP handles a single TCP stream. It sets up a connection to the 'upstream' chat server, relays messages back and forth
// if it finds a BogusCoin address in the message, it will silently rewrite it to match TonysAddress
func (s *Server) HandleTCP(conn net.Conn) {
	// tear down client connection after disconnect
	defer conn.Close()

	// open outbound connection to upstream, fatal if it fails
	oConn, err := net.Dial("tcp", s.Upstream)
	if err != nil {
		s.logger.Err(err).Msg("connecting to upstream")
		return
	}
	defer oConn.Close()
	s.logger.Info().Stringer("upstream", oConn.RemoteAddr()).Stringer("client", conn.RemoteAddr()).Msg("connected")

	// handle the proxying
	go s.pipe(conn, oConn)
	s.pipe(oConn, conn)
	s.logger.Info().Stringer("upstream", oConn.RemoteAddr()).Stringer("client", conn.RemoteAddr()).Msg("disconnected")
}

// pipe reads from connection, rewrites BogusCoin address if present, otherwise writes to the other connection
func (s *Server) pipe(r, w net.Conn) {
	defer func() {
		// we always want to close both sides
		w.Close()
		r.Close()
	}()
	scanner := bufio.NewScanner(r)
	scanner.Split(msgSplitFunc)
	for scanner.Scan() {
		_, err := fmt.Fprintln(w, ReplaceBogusCoins(scanner.Text()))
		if err != nil {
			s.logger.Err(err).Msg("writing")
			return
		}
	}
}

// msgSplitFunc is a bufio.SplitFunc that acts like ScanLines but ignores extra invalid lines on EOF
func msgSplitFunc(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF {
		// if we're already at EOF, we dont want any remaining data
		return 0, nil, nil
	}
	return bufio.ScanLines(data, atEOF)
}
