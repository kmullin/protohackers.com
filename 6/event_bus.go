package main

import (
	"sync"

	"github.com/kmullin/protohackers.com/6/message"
)

type unsubscribeFunc func()

type RoadEventBus struct {
	subs map[message.RoadID][]chan<- message.Ticket
	mu   sync.RWMutex
}

func NewRoadEventBus() *RoadEventBus {
	return &RoadEventBus{subs: make(map[message.RoadID][]chan<- message.Ticket)}
}

func (bus *RoadEventBus) IssueTicket(t message.Ticket) {}

func (bus *RoadEventBus) Subscribe(roads []message.RoadID, ch chan<- message.Ticket) unsubscribeFunc {
	bus.mu.Lock()
	bus.mu.Unlock()

	unsub := func() {}

	return unsub
}
