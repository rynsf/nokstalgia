package driver

import (
	"time"
)

type timer struct {
	id       uint32
	data     uint32
	argc     uint32
	argv     [3]uint32
	last     int64
	interval int64
}

var timers = make([]timer, 0)

func ownTimerStart(id, data, argc uint32, argv [3]uint32, interval int64) {
	now := time.Now().UnixNano()
	t := timer{
		id,
		data,
		argc,
		argv,
		now,
		interval,
	}
	timers = append(timers, t)
}

func ownTimerAbort(id uint32) {
	for i := range timers {
		if id == timers[i].id {
			timers = append(timers[:i], timers[i+1:]...)
		}
	}
}

func timerTick() {
	now := time.Now().UnixNano()
	for i := range timers {
		interval := now - timers[i].last
		if interval >= timers[i].interval {
			timers[i].last = now
			Enq(timers[i].id, timers[i].argc, timers[i].argv)
		}
	}
}
