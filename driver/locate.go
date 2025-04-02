package driver

import (
	"encoding/hex"
	"log"
)

const (
	MEM_RAM = 0x100000
	MEM_MCU = 0x200000
)

var located = map[string]uint32{}
var flash []byte

func InitLocate(f []byte) {
	flash = f
}

func flsGetByte(i uint32) byte {
	return flash[i-MEM_MCU]
}

func flsBound(i uint32) bool {
	if i >= MEM_MCU && i < MEM_MCU+uint32(len(flash)) {
		return true
	}
	return false
}

func x2c(x string) []byte {
	c, err := hex.DecodeString(x)
	if err != nil {
		log.Fatalln("Locate: Invalid pattern")
	}
	return c
}

func locateAuto(id string) uint32 {
	var patt, mask []byte
	switch id {
	case "SEND_MESSAGE":
		patt = x2c("0400B40FB5B0AF04883804C00CC028DCD01B8838")
		mask = x2c("FFFFFFFFFFFFFFFFFFFFFFF0FFFFFFFFFFF0FFFF")
	case "WIN_MARK_DIRTY":
		patt = x2c("B500F7FFFF232800D001F7FFFF3FBD00")
		mask = x2c("FFFFF800F800FFFFFFFFF800F800FFFF")
	default:
		return 0
	}

	loc := find(MEM_MCU, patt, mask)
	if loc != 0 {
		return loc
	}
	return 0
}

func Locate(id string) uint32 {
	handle := func(loc uint32) uint32 {
		located[id] = loc
		return loc
	}
	addr, ok := located[id]
	if ok {
		return addr
	}
	loc := locateAuto(id)
	if loc != 0 {
		return handle(loc)
	}
	log.Println("Locate: object not found")
	return 0
}

func maskedCmp(base uint32, patt, mask []byte) bool {
	if base+uint32(len(patt)) > MEM_MCU+uint32(len(flash)) {
		return false
	}
	for i := range patt {
		if (flsGetByte(base+uint32(i))^patt[i])&mask[i] != 0 {
			return false
		}
	}
	return true
}

func find(base uint32, patt, mask []byte) uint32 {
	i := base
	for flsBound(i) {
		if maskedCmp(i, patt, mask) {
			return i
		}
		i += 1
	}
	return 0
}
