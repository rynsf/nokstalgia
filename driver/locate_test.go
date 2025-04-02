package driver

import (
	"os"
	"testing"
)

var tests = map[string]uint32{
	"SEND_MESSAGE":   0x2F227E,
	"WIN_MARK_DIRTY": 0x265800,
}

func TestLocate(t *testing.T) {
	flash, err := os.ReadFile("./../assets/flash.fls")
	if err != nil {
		t.Errorf("Can not load file")
	}
	InitLocate(flash)
	for key, value := range tests {
		loc := Locate(key)
		if loc != value {
			t.Errorf("Locate %s failed: Expected: %x found: %x", key, value, loc)
		}
	}
}
