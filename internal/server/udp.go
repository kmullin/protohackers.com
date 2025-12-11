package server

import (
	"net"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type UDPHandlerFunc func(net.PacketConn)

func (f UDPHandlerFunc) HandleUDP(conn net.PacketConn) {
	f(conn)
}

type UDPHandler interface {
	HandleUDP(net.PacketConn)
}

func UDP(h UDPHandler) {
	conn, err := net.ListenPacket("udp", viper.GetString("addr"))
	if err != nil {
		log.Fatal().Err(err).Msg("unable to listen")
	}
	log.Info().Stringer("addr", conn.LocalAddr()).Msg("listening")

	h.HandleUDP(conn)
}
