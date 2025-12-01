package main

import (
	"sync"
	"time"

	"github.com/kmullin/protohackers.com/6/message"
	"github.com/rs/zerolog/log"
)

type Observation struct {
	//Plate     string
	Timestamp time.Time
	Mile      uint16
	Limit     uint16
}

type Ticketer struct {
	t  map[uint16]map[string][]Observation
	mu *sync.RWMutex
}

func NewTicketer() *Ticketer {
	return &Ticketer{
		t:  make(map[uint16]map[string][]Observation),
		mu: new(sync.RWMutex),
	}
}

func (t *Ticketer) Observe(plate *message.Plate, camera *message.IAmCamera) {
	t.mu.Lock()
	inner, ok := t.t[camera.Road]
	if !ok {
		inner = make(map[string][]Observation)
		t.t[camera.Road] = inner
	}
	inner[plate.Plate] = append(inner[plate.Plate], Observation{Timestamp: plate.Timestamp, Mile: camera.Mile, Limit: camera.Limit})
	log.Info().Interface("map", t.t).Msg("")
	t.mu.Unlock()
}
