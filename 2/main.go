package main

import (
	"log"

	"git.kpmullin.com/kmullin/protocolhackers.com/2/message"
)

func main() {
	q := message.Query{}
	log.Printf("%v", q)
	i := message.Insert{}
	log.Printf("%v", i)
}
