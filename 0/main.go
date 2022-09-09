// Echo Server
// Implements RFC862
package main

import (
	"io"
	"log"
	"net"
)

func main() {
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Printf("closed %s", conn.RemoteAddr())
	}()
	log.Printf("connected %s", conn.RemoteAddr())
	io.Copy(conn, conn)
}
