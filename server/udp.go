package server

import (
	"net"

	"github.com/rs/zerolog/log"
)

type UDPHandlerFunc func(net.PacketConn)

func (f UDPHandlerFunc) HandleUDP(conn net.PacketConn) {
	f(conn)
}

type UDPHandler interface {
	HandleUDP(net.PacketConn)
}

func UDP(h UDPHandler) {
	addr := GetListenAddr()

	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to listen")
	}
	log.Info().Stringer("addr", conn.LocalAddr()).Msg("listening")

	h.HandleUDP(conn)
}
