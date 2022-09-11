// Echo Server
// Implements RFC862
package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"net"

	"git.kpmullin.com/kmullin/protocolhackers.com/server"
)

func main() {
	server.TCP(server.HandlerFunc(echo))
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
	n, err := io.Copy(w, conn)
	if err != nil {
		log.Printf("err in copy: %v", err)
	}

	enc.Close()
	log.Printf("data received %d: %v", n, buf.String())
}
