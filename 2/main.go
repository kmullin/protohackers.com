package main

import (
	"bufio"
	"log"
	"net"

	"git.kpmullin.com/kmullin/protocolhackers.com/server"
)

func main() {
	server.TCP(handler)
}

func handler(conn net.Conn) {
	defer conn.Close()
	b, err := bufio.NewReader(conn).Peek(1)
	if err != nil {
		log.Printf("err peek: %v", err)
	}
	log.Printf("% x", b)
}
