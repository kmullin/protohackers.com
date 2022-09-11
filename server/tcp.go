package server

import (
	"log"
	"net"
)

type HandlerFunc func(net.Conn)

func (f HandlerFunc) HandleTCP(conn net.Conn) {
	f(conn)
}

type Handler interface {
	HandleTCP(net.Conn)
}

func TCP(h Handler) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("unable to listen: %v", err)
	}
	log.Printf("listening on %v", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept err: %v", err)
			continue
		}
		go h.HandleTCP(conn)
	}
}
