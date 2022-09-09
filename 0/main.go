// Echo Server
// Implements RFC862
package main

import (
	"bytes"
	"encoding/base64"
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
		go echo(conn)
	}
}

func echo(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Printf("closed %s", conn.RemoteAddr())
	}()
	log.Printf("connected %s", conn.RemoteAddr())

	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	w := io.MultiWriter(conn, enc)
	io.Copy(w, conn)

	enc.Close()
	log.Printf("data received: %v", buf.String())
}
