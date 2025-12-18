package main

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/kmullin/protohackers.com/7/message"
	"github.com/stretchr/testify/assert"
)

// Step lets us define the expectations for the server to perform
type Step struct {
	Send   message.Msg
	Expect message.Msg
}

func PacketConnPair(t *testing.T) (*net.UDPConn, *net.UDPConn) {
	t.Helper()

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	server, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	client, err := net.DialUDP("udp", nil, server.LocalAddr().(*net.UDPAddr))
	if err != nil {
		server.Close()
		t.Fatal(err)
	}

	deadline := time.Now().Add(3 * time.Second)
	_ = client.SetDeadline(deadline)
	_ = server.SetDeadline(deadline)

	return client, server
}

func runTranscript(t *testing.T, client *net.UDPConn, steps []Step) {
	t.Helper()

	buf := make([]byte, message.MaxSize)

	for i, step := range steps {
		if step.Send != nil {
			_, err := client.Write(step.Send.Marshal())
			assert.NoErrorf(t, err, "step %d send failed", i)
		}

		if step.Expect != nil {
			n, _, err := client.ReadFromUDP(buf)
			assert.NoErrorf(t, err, "step %d read failed", i)

			msg, err := message.New(buf[:n])
			assert.NoErrorf(t, err, "step %d unmarshal failed", i)

			assert.Equalf(t, step.Expect, msg, "step %d mismatch", i)
		}
	}
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
	client, server := PacketConnPair(t)
	defer client.Close()
	defer server.Close()

	// our client perspective
	steps := []Step{
		{
			Send:   &message.Connect{SessionID: 12345},
			Expect: &message.Ack{SessionID: 12345, Length: 0},
		},
		{
			Send:   &message.Data{SessionID: 12345, Pos: 0, Data: []byte("hello\n")},
			Expect: &message.Ack{SessionID: 12345, Length: 6},
		},
		{
			Send:   &message.Data{SessionID: 12345, Pos: 6, Data: []byte("Hello, world!\n")},
			Expect: &message.Ack{SessionID: 12345, Length: 20},
		},
		{
			Send:   &message.Close{SessionID: 12345},
			Expect: &message.Close{SessionID: 12345},
		},
	}

	var wg sync.WaitGroup
	wg.Add(2)
	// client
	go func() {
		defer wg.Done()

		for _, step := range steps {
			n, err := client.Write(step.Send.Marshal())
			assert.NoError(t, err)
			t.Logf("wrote %v bytes to %v", n, client)

			buf := make([]byte, message.MaxSize)
			n, _, err = client.ReadFrom(buf)
			assert.NoError(t, err)
		}
	}()

	// server
	go func() {
		defer wg.Done()

		buf := make([]byte, message.MaxSize)
		for i := range steps {
			n, addr, err := server.ReadFromUDP(buf)
			t.Logf("server received %d bytes from %v: %q", n, addr, buf[:n])
			assert.NoError(t, err)

			msg, err := message.New(buf[:n])
			assert.NoError(t, err)
			t.Logf("received message: %#v", msg)
			assert.Equal(t, steps[i], msg)
		}
	}()

	wg.Wait()
}
