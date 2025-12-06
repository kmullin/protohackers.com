package main

import (
	"math"
	"sync"

	"github.com/kmullin/protohackers.com/6/message"
	"github.com/rs/zerolog/log"
)

type Ticketer struct {
	obs map[message.RoadID]map[string]Observations
	mu  sync.RWMutex

	subs *RoadEventBus
}

func NewTicketer() *Ticketer {
	return &Ticketer{
		obs:  make(map[message.RoadID]map[string]Observations),
		subs: NewRoadEventBus(),
	}
}

// Observe adds observations to the global state for each plate message from a camera
func (t *Ticketer) Observe(plate *message.Plate, camera *message.IAmCamera) {
	t.mu.Lock()
	defer t.mu.Unlock()

	inner, ok := t.obs[camera.Road]
	if !ok {
		inner = make(map[string]Observations)
		t.obs[camera.Road] = inner
	}

	os := inner[plate.Plate]

	o := Observation{
		Timestamp: plate.Timestamp,
		Mile:      camera.Mile,
		Limit:     camera.Limit,
	}

	os.insertSorted(o)
	inner[plate.Plate] = os

	log.Info().Interface("map", t.obs).Msg("observation")

	t.checkRoad(camera.Road)
}

// checkRoad checks for tickets for a certain road and issues a ticket to the channel
func (t *Ticketer) checkRoad(road message.RoadID) {
	// mutex already locked
	for plate, observations := range t.obs[road] {
		if len(observations) < 2 {
			// we need at least 2 observations to determine speed
			continue
		}

		for i := 0; i < len(observations)-1; i++ {
			a := observations[i]
			b := observations[i+1]

			distance := math.Abs(float64(a.Mile) - float64(b.Mile))
			duration := b.Timestamp.Sub(a.Timestamp).Hours()
			speed := distance / duration

			log.Info().
				Interface("a", a).
				Interface("b", b).
				Float64("speed", speed).
				Float64("distance", distance).
				Float64("duration", duration).
				Str("plate", plate).
				Uint16("road", uint16(road)).
				Uint16("limit", a.Limit).
				Msg("get speed")

			// always required to ticket a car that is exceeding the speed limit by 0.5 mph or more
			if speed >= (float64(a.Limit) + 0.5) {
				ticket := message.Ticket{
					Plate:      plate,
					Road:       road,
					Mile1:      a.Mile,
					Timestamp1: a.Timestamp,
					Mile2:      b.Mile,
					Timestamp2: b.Timestamp,
					Speed:      uint16(speed),
				}
				log.Info().
					Interface("ticket", ticket).
					Msg("issuing ticket")

				t.subs.IssueTicket(ticket)
				// XXX: need to delete observations
			}
		}
	}
}
