package emulatedispatcher

import (
	dr "github.com/rynsf/nokstalgia/driver"
)

func (s *CpuState) subroutineIsRunning() bool {
	return s.register[pc] != 0
}

func InitCpu(r [16]uint32, v, c, z, n bool, ram []byte, ramBase, ramLen uint32, dynamicRam []byte, dynRamBase, dynRamLen uint32, flash []byte, flashBase, flashLen uint32) CpuState {
	specialFunc[dr.Locate("MALLOC")] = malloc
	dr.InitDynamicMem(0x4, 0xe000)
	specialFunc[dr.Locate("MARK_DIRTY")] = mark_dirty
	specialFunc[dr.Locate("WIN_PRINT")] = win_print
	specialFunc[dr.Locate("OWN_TIMER_START")] = own_timer_start
	specialFunc[dr.Locate("OWN_TIMER_ABORT")] = own_timer_abort
	specialFunc[dr.Locate("SEND_MESSAGE")] = send_message
	specialFunc[dr.Locate("OWN_GET_CONFIG")] = own_get_config
	specialFunc[dr.Locate("OWN_GET_FONT")] = own_get_font
	return CpuState{
		register:     r,
		sr:           cpsr{v, c, z, n},
		loc:          0,
		ram:          ram,
		ramBaseAdr:   ramBase,
		ramLen:       ramLen,
		dynamicRam:   dynamicRam,
		dynRamBase:   dynRamBase,
		dynRamLen:    dynRamLen,
		flash:        flash,
		flashBaseAdr: flashBase,
		flashLen:     flashLen,
	}
}

func (s *CpuState) RunSubroutine() {
	for s.subroutineIsRunning() {
		if s.specialFuncHandler() {
			continue
		}
		instruction := uint16(s.read16(s.register[pc]))
		if Debug {
			s.step()
		}
		s.loc = s.register[pc]
		s.register[pc] += 2
		s.execInstruction(instruction)
	}
}
