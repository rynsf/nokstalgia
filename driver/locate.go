package driver

import (
	"encoding/hex"
	"log"
)

const (
	MEM_RAM = 0x100000
	MEM_MCU = 0x200000
)

// predefined addresses for now.
// TODO: write locate functions for these.
var located = map[string]uint32{
	"SEND_MESSAGE":                  0x2F227E,
	"BLOCK_ALLOC":                   0x2ACD14,
	"BLOCK_DEALLOC":                 0x2ACCA4,
	"BLOCK_ALLOC_NOWAIT":            0x2ACD94,
	"MALLOC":                        0x2F9166,
	"FREE":                          0x2F9300,
	"DEV_FUNC_TRACE":                0x307C04,
	"TI_ID_SEND":                    0x301E4C,
	"OS_CONDITIONAL_INT_ENABLE":     0x2FA834,
	"OS_CONDITIONAL_INT_DISABLE":    0x2FA812,
	"DEV_DISP_REFRESH":              0x2F9A0E,
	"SCREEN_BUFFER":                 0x140012,
	"SCREEN_HEIGHT":                 0x41,
	"SCREEN_WIDTH":                  0x60,
	"OWN_TIMER_START":               0x3D43F0,
	"LOAD_GLOBAL_SETTINGS_VALUE":    0x2CECE4,
	"READ_DIRECTORY_FILE":           0x2F5D82,
	"ENGINE_INITILIZE_FILE":         0x2CF114,
	"TRANSFER_DOWNLOADED_DATA_INFO": 0x2CE72C,
	"WRITE_DIRECTORY_FILE":          0x2CE72C,
	"FREE_DIRECTORY_FILE":           0x2F610C,
	"CREATE_MENU":                   0x307546,
	"LOAD_SAVED_GAME":               0x2E5DB8,
	"TRANSLATE_UCS2":                0x2BF8FE,
	"OS_TIMER_STOP":                 0x2ABA4C,
	"OS_TIMER_START":                0x2ABBD0,
	"IND_CALL":                      0x308880,
	"MSG_ID":                        0x10D9FA,
	"MSG_ARGC":                      0x10D9F4,
	"MSG_ARGV":                      0x10DA0C,
	"DEV_KEY_GAME_MODE_ENABLE":      0x300FC8,
	"TONE_CLASS_DISABLE":            0x301CF6,
	"GAME_LOAD_HIGHSCORE":           0x2F4D98,
	"ENGINE_LOAD_SETTINGS_VALUE":    0x2CEBA0,
	"RAM_INIT_BLOCK":                0x3088D0,
}

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
	log.Printf("Locate: object %s not found\n", id)
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
