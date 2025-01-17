package emulatedispatcher

func (s *CpuState) addOffsetSP(instruction uint16) {
	offset := getBitsRange(instruction, 0, 6)
	offset *= 4
	op := getBitsRange(instruction, 7, 7)

	switch op {
	case 0:
		s.register[sp] += uint32(offset)
	case 1:
		s.register[sp] -= uint32(offset)
	}
}
