package main

import (
	"encoding/binary"
	"io"
	"log"
	"math"
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

			mean := findMean(inserts, m.MinTime, m.MaxTime)
			err := binary.Write(conn, message.ByteOrder, mean)
			if err != nil {
				log.Printf("err sending response: %v", err)
			}
		default:
			log.Printf("not implemented yet: %+v", m)
		}
	}
}

func findMean(inserts []message.Insert, min, max time.Time) int32 {
	var count, result int32
	if min.Before(max) {
		for _, i := range inserts {
			if i.Timestamp.Before(min) || i.Timestamp.After(max) {
				continue
			}
			// closed interval, we're equal or between
			result += i.Price
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return int32(math.Round(float64(result) / float64(count)))
}
