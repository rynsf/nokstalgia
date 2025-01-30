package emulatedispatcher

func (s *CpuState) subroutineIsRunning() bool {
	return s.register[pc] != 0
}

func InitCpu(r [16]uint32, v, c, z, n bool, ram []byte, ramBase, ramLen uint32, dynamicRam []byte, dynRamBase, dynRamLen uint32, flash []byte, flashBase, flashLen uint32) CpuState {
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
		instruction := uint16(s.read16(s.register[pc]))
		if Debug {
			s.step()
		}
		s.loc = s.register[pc]
		s.register[pc] += 2
		s.execInstruction(instruction)
	}
}
