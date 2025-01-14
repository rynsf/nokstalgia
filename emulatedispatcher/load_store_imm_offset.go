package emulatedispatcher

func (s *CpuState) LoadStoreImmOffset(instruction uint16) {
	rd := getBitsRange(instruction, 0, 2)
	rb := getBitsRange(instruction, 3, 5)
	offset := getBitsRange(instruction, 6, 10)
	op := getBitsRange(instruction, 11, 12)

	switch op {
	case 0: // STR
		offset *= 4
		s.write32(s.register[rb]+uint32(offset), s.register[rd])
	case 1: // LDR
		offset *= 4
		s.register[rd] = s.read32(s.register[rb] + uint32(offset))
	case 2: // STRB
		s.write8(s.register[rb]+uint32(offset), byte(s.register[rd]))
	case 3: // LDRB
		s.register[rd] = uint32(s.read8(s.register[rb] + uint32(offset)))
	}
}
