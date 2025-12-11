package main

import (
	"sync"

	"github.com/kmullin/protohackers.com/6/message"
	"github.com/rs/zerolog"
)

const ticketCSize = 8192

type RoadEventBus struct {
	subs map[message.RoadID]chan message.Ticket
	mu   sync.RWMutex

	log zerolog.Logger
}

func NewRoadEventBus(logger zerolog.Logger) *RoadEventBus {
	return &RoadEventBus{
		subs: make(map[message.RoadID]chan message.Ticket),
		log:  logger.With().Str("component", "eventbus").Logger(),
	}
}

func (b *RoadEventBus) IssueTicket(t message.Ticket) {
	b.mu.RLock()
	ch := b.subs[t.Road]
	b.mu.RUnlock()

	ch <- t
	b.log.Info().
		Interface("ticket", t).
		Msg("issuing ticket")
}

// NewRoad creates fresh channel for the road to write tickets to
func (b *RoadEventBus) NewRoad(road message.RoadID) <-chan message.Ticket {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch, ok := b.subs[road]
	if !ok || ch == nil {
		ch = make(chan message.Ticket, ticketCSize)
		b.subs[road] = ch
	}
	return ch
}

func (b *RoadEventBus) Subscribe(done <-chan struct{}, roads []message.RoadID) <-chan message.Ticket {
	var chs []<-chan message.Ticket
	out := make(chan message.Ticket, ticketCSize*len(roads))

	// group all road channels into one
	for _, road := range roads {
		chs = append(chs, b.NewRoad(road))
	}

	for _, ch := range chs {
		go func(c <-chan message.Ticket) {
			for {
				select {
				case <-done:
					return
				case v := <-c:
					out <- v
				}
			}
		}(ch)
	}

	go func() {
		<-done
		close(out)
	}()

	return out
}
