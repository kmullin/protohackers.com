package test

import (
	"net"
	"sync"
	"testing"
	"time"
)

// taken from https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=37
func Conn(t *testing.T) (client, server net.Conn) {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer ln.Close()

		var err error // otherwise a data race trying to set err
		server, err = ln.Accept()
		if err != nil {
			t.Fatal(err)
		}
	}()

	client, err = net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()
	// configure some upper limit deadline for tests
	now := time.Now()
	for _, conn := range []net.Conn{client, server} {
		conn.SetDeadline(now.Add(3 * time.Second))
	}
	return client, server
}
