package main

import (
	"fmt"
	"reflect"
	"testing"

	"git.kpmullin.com/kmullin/protocolhackers.com/test"
)

func TestEcho(t *testing.T) {
	cases := [][]byte{
		[]byte("foobar"),
		[]byte("oijaiojdoiuas0dajs9djpaskd 09aj09 j09 j09jsdoj odjfg j1lks;ldfk;lsjgih98"),
		[]byte("\nbaz"),
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			client, server := test.Conn(t)
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
