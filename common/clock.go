package common

import (
	"fmt"
	"sync"
)

type VectorClock struct {
	localPid int
	state    []int
}

// a singleton
var clock *VectorClock
var vcOnce sync.Once

func GetLocalVectorClock() *VectorClock {
	return clock
}

// creates the singleton
func NewLocalVectorClock(peers, pid int) *VectorClock {
	vcOnce.Do(func() {
		clock = new(VectorClock)
		clock.state = make([]int, peers)
		clock.localPid = pid
	})

	return clock
}

func NewVectorClock(other []int) *VectorClock {
	vc := new(VectorClock)
	vc.state = other
	return vc
}

func (me *VectorClock) IncrementClock() {
	me.state[me.localPid]++
}

func (me *VectorClock) UpdateClock(other *VectorClock) {
	for i, v := range me.state {
		if v < other.state[i] {
			me.state[i] = other.state[i]
		}
	}
}

// true if me happened before other
func (me *VectorClock) HappenedBefore(other *VectorClock) bool {
	for i, v := range me.state {
		if v > other.state[i] {
			return false
		}
	}
	return true
}

func (me *VectorClock) HappenedAfter(other *VectorClock) bool {
	for i, v := range me.state {
		if v < other.state[i] {
			return false
		}
	}
	return true
}

func (me *VectorClock) Independent(other *VectorClock) bool {
	return !me.HappenedBefore(other) && !me.HappenedAfter(other)
}

func (me *VectorClock) String() string {
	return fmt.Sprintf("%v", me.state)
}