package emulatedispatcher

func align4(val uint32) uint32 {
	return val & ^uint32(0b11)
}

func (s *CpuState) loadAddress(instruction uint16) {
	offset := getBitsRange(instruction, 0, 7)
	offset *= 4
	rd := getBitsRange(instruction, 8, 10)
	op := getBitsRange(instruction, 11, 11)

	switch op {
	case 0: // PC
		s.register[rd] = (align4(s.loc + 4)) + uint32(offset)
	case 1: // SP
		s.register[rd] = s.register[sp] + uint32(offset)
	}
}
