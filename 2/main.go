package main

import (
	"io"
	"log"
	"net"

	"git.kpmullin.com/kmullin/protocolhackers.com/2/message"
	"git.kpmullin.com/kmullin/protocolhackers.com/server"
)

func main() {
	server.TCP(server.HandlerFunc(handler))
}

func handler(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Printf("closed: %v", conn.RemoteAddr())
	}()

	for {
		m, err := message.New(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("msg err: %v", err)
			return
		}
		log.Printf("%+v", m)
	}
}
