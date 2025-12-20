package main

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/kmullin/protohackers.com/7/message"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// Step lets us define the expectations for the server to perform
type Step struct {
	Send   message.Msg
	Expect []message.Msg
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

		for ii, expect := range step.Expect {
			n, _, err := client.ReadFromUDP(buf)
			assert.NoErrorf(t, err, "step %d expect %v read failed", i, ii)

			msg, err := message.New(buf[:n])
			assert.NoErrorf(t, err, "step %d expect %v unmarshal failed", i, ii)

			assert.Equalf(t, expect, msg, "step %d expect %v mismatch", i, ii)
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
			Send: &message.Connect{SessionID: 12345},
			Expect: []message.Msg{
				&message.Ack{SessionID: 12345, Length: 0},
			},
		},
		{
			Send: &message.Data{SessionID: 12345, Pos: 0, Data: []byte("hello\n")},
			Expect: []message.Msg{
				&message.Ack{SessionID: 12345, Length: 6},
				&message.Data{SessionID: 12345, Pos: 0, Data: []byte("olleh\n")},
			},
		},
		{
			Send: &message.Data{SessionID: 12345, Pos: 6, Data: []byte("Hello, world!\n")},
			Expect: []message.Msg{
				&message.Ack{SessionID: 12345, Length: 20},
				&message.Data{SessionID: 12345, Pos: 6, Data: []byte("!dlrow ,olleH\n")},
			},
		},
		{
			Send: &message.Close{SessionID: 12345},
			Expect: []message.Msg{
				&message.Close{SessionID: 12345},
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start our server
	go func() {
		s := NewServer(ctx)
		s.log = testLogger(t)
		s.HandleUDP(server)
	}()

	runTranscript(t, client, steps)
}

func testLogger(t *testing.T) zerolog.Logger {
	t.Helper()
	cw := zerolog.ConsoleWriter{}
	zerolog.ConsoleTestWriter(t)(&cw)
	return zerolog.New(cw)
}
