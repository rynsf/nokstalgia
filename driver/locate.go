package driver

import (
	"log"
)

var located = map[string]uint32{
	"MALLOC":          0x2551F4,
	"MARK_DIRTY":      0x2ADB0C,
	"WIN_PRINT":       0x234B7A,
	"MSG_ID":          0x10A4EA,
	"MSG_ARGC":        0x10A4E6,
	"MSG_ARGV":        0x10A4EC,
	"OWN_TIMER_START": 0x2AD494,
	"OWN_TIMER_ABORT": 0x2AD3F8,
	"SEND_MESSAGE":    0x28D21A,
}

// First and a very basic implementation of locate function.
// locate takes in name of an object and returns its memory address.
func Locate(id string) uint32 {
	addr, ok := located[id]
	if ok {
		return addr
	}
	log.Println("Locate: object not found")
	return 0
}
