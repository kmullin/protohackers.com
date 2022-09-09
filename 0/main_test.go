package main

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"
)

func testServer(t *testing.T) (client, server net.Conn) {
	t.Helper()

	var err error
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer ln.Close()
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
	return client, server
}

func TestEcho(t *testing.T) {
	cases := [][]byte{
		[]byte("foobar"),
		[]byte("oijaiojdoiuas0dajs9djpaskd 09aj09 j09 j09jsdoj odjfg j1lks;ldfk;lsjgih98"),
		[]byte("\nbaz"),
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			client, server := testServer(t)
			go echo(server)

			n, err := client.Write(tc)
			if err != nil {
				t.Fatal(err)
			}
			if n != len(tc) {
				t.Fatalf("len written is not expected, got %v wanted %v", n, len(tc))
			}

			b := make([]byte, len(tc))
			n, err = client.Read(b)
			if err != nil {
				t.Fatal(err)
			}
			if n != len(tc) {
				t.Fatalf("len written is not expected, got %v wanted %v", n, len(tc))
			}

			if !reflect.DeepEqual(tc, b) {
				t.Fatalf("input and output differ, got %v, wanted %v", b, tc)
			}

			client.Close()
		})
	}
}
