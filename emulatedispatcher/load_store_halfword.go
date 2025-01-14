package emulatedispatcher

func (s *CpuState) loadStoreHalfword(instruction uint16) {
	rd := getBitsRange(instruction, 0, 2)
	rb := getBitsRange(instruction, 3, 5)
	offset := getBitsRange(instruction, 6, 10)
	offset *= 2
	lflag := getBitsRange(instruction, 11, 11)

	switch lflag {
	case 0: // STRH
		s.write16(s.register[rb]+uint32(offset), uint16(s.register[rd]))
	case 1: // LDRH
		s.register[rd] = uint32(s.read16(s.register[rb] + uint32(offset)))
	}
}
