package emulatedispatcher

func (s *CpuState) spRelLoadStore(instruction uint16) {
	offset := getBitsRange(instruction, 0, 7)
	offset *= 4
	rd := getBitsRange(instruction, 8, 10)
	op := getBitsRange(instruction, 11, 11)
	sp := s.register[sp]

	switch op {
	case 0: // STR
		s.write32(sp+uint32(offset), s.register[rd])
	case 1: // LDR
		s.register[rd] = s.read32(sp + uint32(offset))
	}
}
