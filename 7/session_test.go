package main

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/kmullin/protohackers.com/7/message"
	"github.com/stretchr/testify/assert"
)

func PacketConn(t *testing.T) (client, server net.PacketConn) {
	t.Helper()

	// Server listens
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	server, err = net.ListenUDP("udp", serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	// Client dials server
	client, err = net.DialUDP("udp", nil, server.LocalAddr().(*net.UDPAddr))
	if err != nil {
		server.Close()
		t.Fatal(err)
	}

	// Set deadlines so tests don't hang
	deadline := time.Now().Add(3 * time.Second)
	client.SetDeadline(deadline)
	server.SetDeadline(deadline)

	return client, server
}

/*
<-- /connect/12345/
--> /ack/12345/0/
<-- /data/12345/0/hello\n/
--> /ack/12345/6/
--> /data/12345/0/olleh\n/
<-- /ack/12345/6/
<-- /data/12345/6/Hello, world!\n/
--> /ack/12345/20/
--> /data/12345/6/!dlrow ,olleH\n/
<-- /ack/12345/20/
<-- /close/12345/
--> /close/12345/
*/
func TestSessionExample(t *testing.T) {
	client, server := PacketConn(t)
	defer client.Close()
	defer server.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		msgs := []message.Msg{
			&message.Connect{SessionID: 12345},
			&message.Data{SessionID: 12345, Pos: 0, Data: []byte("hello\n")},
			&message.Data{SessionID: 12345, Pos: 6, Data: []byte("Hello, world!\n")},
			&message.Close{SessionID: 12345},
		}

		for _, msg := range msgs {
			n, err := client.Write(msg.Marshal())
			assert.NoError(t, err)
			t.Logf("wrote %v bytes to %v", n, conn.LocalAddr())
		}
	}()

	t.Logf("conn: %+v", conn)

	// buf := make([]byte, message.MaxSize)
	// n, addr, err := conn.ReadFrom(buf)
	// assert.NoError(t, err)

	// t.Logf("read %v bytes from %v", n, addr)
	wg.Wait()
}
