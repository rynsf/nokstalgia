package emulatedispatcher

import (
	driver "github.com/rynsf/nokstalgia/driver"
)

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
	srcPtr := uint32(0x101418)
	src := s.read32(srcPtr)
	dst := uint32(0x107604)
	for i := 0; i < (84 * 6); i++ {
		b := s.read8(src + uint32(i))
		s.write(dst+uint32(i), b)
	}
}

// TODO: implement blink buffer
func (s *CpuState) SendToLcd(screen [][]int) [][]int {
	screenBufferBase := 0x107604
	for y := 0; y < 48; y++ {
		for x := 0; x < 84; x++ {
			yByteAdr := y / 8
			yBitAdr := y % 8
			pixelByte := s.read(uint32(screenBufferBase + (yByteAdr * 84) + x))
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
