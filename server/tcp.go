package server

import (
	"log"
	"net"
)

type HandleFunc func(net.Conn)

func TCP(handler HandleFunc) {
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
		go handler(conn)
	}
}
