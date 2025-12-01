package main

import (
	"math"
	"sort"
	"sync"
	"time"

	"github.com/kmullin/protohackers.com/6/message"
	"github.com/rs/zerolog/log"
)

type Observation struct {
	Timestamp time.Time
	Mile      uint16
	Limit     uint16
}

type Observations []Observation

type Ticketer struct {
	t  map[uint16]map[string]Observations
	mu *sync.RWMutex
}

// insertSorted inserts into the slice in a sorted manner
func (o *Observations) insertSorted(obs Observation) {
	s := *o

	// find insertion index using binary search
	i := sort.Search(len(s), func(i int) bool {
		return s[i].Timestamp.After(obs.Timestamp) || s[i].Timestamp.Equal(obs.Timestamp)
	})

	// expand slice by 1
	s = append(s, Observation{})

	// shift elements to the right
	copy(s[i+1:], s[i:])

	// insert the new element
	s[i] = obs

	*o = s
}

func NewTicketer() *Ticketer {
	return &Ticketer{
		t:  make(map[uint16]map[string]Observations),
		mu: new(sync.RWMutex),
	}
}

func (t *Ticketer) Observe(plate *message.Plate, camera *message.IAmCamera) {
	t.mu.Lock()
	inner, ok := t.t[camera.Road]
	if !ok {
		inner = make(map[string]Observations)
		t.t[camera.Road] = inner
	}

	os := inner[plate.Plate]

	o := Observation{
		Timestamp: plate.Timestamp,
		Mile:      camera.Mile,
		Limit:     camera.Limit,
	}

	os.insertSorted(o)
	inner[plate.Plate] = os

	log.Info().Interface("map", t.t).Msg("")
	t.mu.Unlock()
}

func (t *Ticketer) Check(roads []uint16) *message.Ticket {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, road := range roads {
		for plate, observations := range t.t[road] {
			if len(observations) < 2 {
				continue
			}

			for i := 0; i < len(observations)-1; i++ {
				a := observations[i]
				b := observations[i+1]

				distance := math.Abs(float64(a.Mile) - float64(b.Mile))
				duration := b.Timestamp.Sub(a.Timestamp).Hours()
				speed := uint16(distance / duration * 100)

				log.Info().
					Interface("a", a).
					Interface("b", b).
					Uint16("speed", speed).
					Float64("distance", distance).
					Float64("duration", duration).
					Str("plate", plate).
					Uint16("road", road).
					Uint16("limit", a.Limit).
					Msg("getSpeed")

			}
		}
	}

	return nil
}
