package emulatedispatcher

import (
	driver "github.com/rynsf/nokstalgia/driver"
)

const TPS = (1000 / 7.96875) // ticks per second

var specialFunc = make(map[uint32]func(*CpuState))

/*
Special functions are low level operating system functions.
As this emulator emulates both hardware and the operating
system of the phone, it handles very low level function
separately. This function checks if the pc points to one such
special function using a pre-initalize map and calls that
function.
*/
func (s *CpuState) specialFuncHandler() bool {
	f, ok := specialFunc[s.register[pc]]
	if ok {
		f(s)
		s.register[pc] = s.register[lr]
		return true
	}
	return false
}

// malloc is a wrapper to driver.malloc. It takes in the parameter for malloc from the register and set registers for output.
func malloc(s *CpuState) {
	size := s.register[0]
	ptr := driver.Malloc(size)
	ptr += s.dynRamBase
	s.register[0] = ptr
}

// free is a wrapper for driver.free
func free(s *CpuState) {
	ptr := s.register[0]
	ptr -= s.dynRamBase
	driver.Free(ptr)
}

func (s *CpuState) UpdateScreen() {
	s.ownDrawingRoutine()
}

// nokix's own window drawing routine
func (s *CpuState) ownDrawingRoutine() {
	srcPtr := uint32(0x101424)
	src := s.read32(srcPtr) + 8
	dst := uint32(0x107604)
	for i := 0; i < (84 * 6); i++ {
		b := s.read8(src + uint32(i))
		s.write(dst+uint32(i), b)
	}
}

// TODO: implement blink buffer
func (s *CpuState) SendToLcd(screen [][]int) [][]int {
	screenBufferBase := int(driver.Locate("SCREEN_BUFFER"))
	h := int(driver.Locate("SCREEN_HEIGHT"))
	w := int(driver.Locate("SCREEN_WIDTH"))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			yByteAdr := y / 8
			yBitAdr := y % 8
			pixelByte := s.read(uint32(screenBufferBase + (yByteAdr * w) + x))
			pixelBit := pixelByte & (1 << yBitAdr)
			if pixelBit == 0 {
				screen[y][x] = 0
			} else {
				screen[y][x] = 1
			}
		}
	}
	return screen
}

// TODO: write a more complete implementation of windowing system later. But it is not need for current implementation of update_screen.
func win_print(s *CpuState) {
}

// TODO: write a more complete implementation of windowing system later. But it is not need for current implementation of update_screen.
func mark_dirty(s *CpuState) {
}

func (s *CpuState) SetMessage(id, argc uint32, argv [3]uint32) {
	addrId := driver.Locate("MSG_ID")
	addrArgc := driver.Locate("MSG_ARGC")
	addrArgv := driver.Locate("MSG_ARGV")
	s.write16(addrId, uint16(id))
	s.write8(addrArgc, uint8(argc))
	for i := 0; i < 3; i++ {
		s.write32(addrArgv+uint32(i*4), argv[i])
	}
}

func os_timer_start(s *CpuState) {
	id := s.register[0]
	t := s.register[1]
	msg := uint32(0x5AF)        // games only use this timer message, TODO: implement loading data for other timers
	interval := t * (1e9 / TPS) // convert ticks to nanoseconds
	driver.TimerStart(id, msg, 0, [3]uint32{}, int64(interval))
}

func os_timer_stop(s *CpuState) {
	id := s.register[0]
	driver.TimerStop(id)
}

func send_message(s *CpuState) {
}

func own_get_config(s *CpuState) {
	switch s.register[0] {
	case 0x0: // level config
		s.register[0] = 0x7
	case 0x1: // maze config
		s.register[0] = 0x0
	}
}

func own_get_font(s *CpuState) {
	s.register[0] = 0x2B04EC
	s.register[1] = 0x0
}

func block_alloc(s *CpuState) {
	a1 := int32(s.register[0]) // cast to signed int32 to get arithmetic shift
	a1 = (a1 + 0x23) >> 4
	size := uint32(a1*16) - 4
	ptr := driver.Malloc(size)
	ptr += s.dynRamBase
	s.register[0] = ptr
}

func block_dealloc(s *CpuState) {
	ptr := s.register[0]
	ptr -= s.dynRamBase
	driver.Free(ptr)
}

func block_alloc_nowait(s *CpuState) {
	block_alloc(s)
}

func doNothing(s *CpuState) {
}

func loadGlobalSettingsValue(s *CpuState) {
	s.register[0] = 0x1
}

func readDirectoryFile(s *CpuState) {
	s.register[0] = 0x1
}

func gameLoadHighscore(s *CpuState) {
	// allocate space
	ptr := driver.Malloc(28)
	ptr += s.dynRamBase
	// fill zeros
	for i := uint32(0); i < 28; i += 4 {
		s.write32(ptr+i, 0x0) // we don't have pmm, just returning 0s
	}
	// write pointer
	dstPtr := s.register[1]
	s.write32(dstPtr, ptr)
	// return 1
	s.register[0] = 0x1
}

func engineLoadSettingsValue(s *CpuState) {
	s.register[0] = 0
}
