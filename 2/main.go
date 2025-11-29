// Means to an End
package main

import (
	"log"
	"math"
	"net"
	"time"

	"github.com/kmullin/protohackers.com/2/message"
	"github.com/kmullin/protohackers.com/server"
)

func main() {
	server.TCP(server.HandlerFunc(handler))
}

type insertCache map[time.Time]int32

func (c insertCache) Mean(min, max time.Time) int32 {
	var count float64
	var result int64
	if !min.After(max) {
		for t, p := range c {
			if t.Before(min) || t.After(max) {
				continue
			}
			// closed interval, we're equal or between
			result += int64(p)
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return int32(math.Round(float64(result) / count))
}

func handler(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Printf("closed: %v", conn.RemoteAddr())
	}()

	inserts := make(insertCache)
	for {
		i, err := message.New(conn)
		if err != nil {
			log.Printf("msg err: %v", err)
			return
		}
		switch m := i.(type) {
		case message.Insert:
			inserts[m.Timestamp] = m.Price
		case message.Query:
			log.Printf("query: %v", m)
			err := message.Write(conn, inserts.Mean(m.MinTime, m.MaxTime))
			if err != nil {
				log.Printf("err sending response: %v", err)
			}
		default:
			log.Printf("not implemented yet: %+v", m)
		}
	}
}
