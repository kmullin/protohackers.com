package main

import (
	"io"
	"log"
	"net"
	"time"

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

	var inserts []message.Insert

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
			inserts = append(inserts, m)
		case message.Query:
			log.Printf("query: %v", m)
		default:
			log.Printf("not implemented yet: %+v", m)
		}
	}

	log.Printf("inserts: %v", inserts)
}

func findMean(inserts []message.Insert, min, max time.Time) (result int32) {
	var count int32
	for _, i := range inserts {
		if i.Timestamp.After(min) && i.Timestamp.Before(max) {
			result += i.Price
			count++
		}
	}
	return result / count
}
