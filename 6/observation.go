package main

import (
	"sort"
	"time"
)

type Observation struct {
	Timestamp time.Time
	Mile      uint16
	Limit     uint16
}

type Observations []Observation

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
