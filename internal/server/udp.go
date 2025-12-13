package server

import (
	"fmt"
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

func UDP(h UDPHandler) error {
	conn, err := net.ListenPacket("udp", viper.GetString("addr"))
	if err != nil {
		return fmt.Errorf("unable to listen: %w", err)
	}
	log.Info().Stringer("addr", conn.LocalAddr()).Msg("listening")
	defer conn.Close()

	h.HandleUDP(conn)
	return nil
}
