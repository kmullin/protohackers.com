package main

import (
	"math"
	"sync"
	"time"

	"github.com/kmullin/protohackers.com/6/message"
	"github.com/rs/zerolog"
)

type TicketRecorder struct {
	record map[string]map[time.Time]struct{}
	mu     sync.RWMutex
}

func NewTicketRecorder() *TicketRecorder {
	return &TicketRecorder{
		record: make(map[string]map[time.Time]struct{}),
	}
}

func (tr *TicketRecorder) Record(t message.Ticket) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if _, ok := tr.record[t.Plate]; !ok {
		tr.record[t.Plate] = make(map[time.Time]struct{})
	}

	for _, ts := range []time.Time{t.Timestamp1, t.Timestamp2} {
		// store the day at midnight
		tr.record[t.Plate][midnight(ts)] = struct{}{}
	}
}

func (tr *TicketRecorder) ExistingTicket(t message.Ticket) bool {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	for _, ts := range []time.Time{t.Timestamp1, t.Timestamp2} {
		_, ok := tr.record[t.Plate][midnight(ts)]
		if ok {
			return true
		}
	}
	return false
}

type Ticketer struct {
	obs map[message.RoadID]map[string]Observations
	mu  sync.RWMutex

	record *TicketRecorder
	subs   *RoadEventBus

	logger zerolog.Logger
}

func NewTicketer(logger zerolog.Logger) *Ticketer {
	return &Ticketer{
		obs:    make(map[message.RoadID]map[string]Observations),
		subs:   NewRoadEventBus(logger),
		record: NewTicketRecorder(),
		logger: logger.With().Str("component", "ticketer").Logger(),
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

	t.logger.Info().Str("plate", plate.Plate).Uint16("road", uint16(camera.Road)).
		Interface("observations", os).Msg("new observation")

	if len(os) >= 2 {
		t.checkRoad(plate.Plate, camera.Road)
	}
}

// checkRoad checks for tickets for a certain road and issues a ticket to the channel
func (t *Ticketer) checkRoad(plate string, road message.RoadID) {
	// mutex already locked
	observations := t.obs[road][plate]
	for i := 0; i < len(observations)-1; i++ {
		a := observations[i]
		b := observations[i+1]

		distance := math.Abs(float64(a.Mile) - float64(b.Mile))
		duration := b.Timestamp.Sub(a.Timestamp).Hours()
		speed := distance / duration

		t.logger.Info().
			Interface("a", a).
			Interface("b", b).
			Float64("speed", speed).
			Float64("distance", distance).
			Float64("duration", duration).
			Str("plate", plate).
			Uint16("road", uint16(road)).
			Uint16("limit", a.Limit).
			Msg("check speed")

		// always required to ticket a car that is exceeding the speed limit by 0.5 mph or more
		if speed >= (float64(a.Limit) + 0.5) {
			ticket := message.Ticket{
				Plate:      plate,
				Road:       road,
				Mile1:      a.Mile,
				Timestamp1: a.Timestamp,
				Mile2:      b.Mile,
				Timestamp2: b.Timestamp,
				Speed:      uint16(math.Round(speed)),
			}

			if t.record.ExistingTicket(ticket) {
				t.logger.Info().Interface("ticket", ticket).
					Msg("already issued for day")
				continue
			}

			t.subs.IssueTicket(ticket)
			t.record.Record(ticket)

			// since ticket issued, delete ticketed observations
			observations = append(observations[:i], observations[i+2:]...)
			t.obs[road][plate] = observations
			i--
			t.logger.Info().
				Interface("new observations", observations).
				Msg("deleted 2 observations")
		}
	}
}

// IsSameDay determines if two timestamps are the same day
func IsSameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()

	return ay == by && am == bm && ad == bd
}

func midnight(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
