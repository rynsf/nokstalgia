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
