package emulatedispatcher

import (
	driver "github.com/rynsf/nokstalgia/driver"
)

// malloc is a wrapper to driver.malloc. It takes in the parameter for malloc from the register and set registers for output.
func (s *CpuState) malloc() {
	size := s.register[0]
	ptr := driver.Malloc(size)
	ptr += s.dynRamBase
	s.register[0] = ptr
}

// free is a wrapper for driver.free
func (s *CpuState) free() {
	ptr := s.register[0]
	ptr -= s.dynRamBase
	driver.Free(ptr)
}
