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
		i, err := message.New(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("msg err: %v", err)
			return
		}
		switch m := i.(type) {
		case message.Insert:
			log.Printf("insert: %v", m)
		case message.Query:
			log.Printf("query: %v", m)
		default:
			log.Printf("not implemented yet: %+v", m)
		}
	}
}
