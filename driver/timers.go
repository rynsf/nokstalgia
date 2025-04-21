package driver

import (
	"time"
)

type timer struct {
	id   uint32
	msg  uint32
	argc uint32
	argv [3]uint32
	end  int64
}

var timers = make(map[uint32]timer)

func TimerStart(id, msg, argc uint32, argv [3]uint32, interval int64) {
	end := time.Now().UnixNano() + interval
	t := timer{
		id,
		msg,
		argc,
		argv,
		end,
	}
	timers[t.id] = t
}

func TimerStop(id uint32) {
	delete(timers, id)
}

func TimerTick() {
	now := time.Now().UnixNano()
	for i := range timers {
		if now >= timers[i].end {
			Enq(timers[i].msg, timers[i].argc, timers[i].argv)
			delete(timers, i)
		}
	}
}
