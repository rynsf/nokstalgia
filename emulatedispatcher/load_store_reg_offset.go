package emulatedispatcher

func (s *CpuState) loadStoreRegOffset(instruction uint16) {
	ro := getBitsRange(instruction, 6, 8)
	rb := getBitsRange(instruction, 3, 5)
	rd := getBitsRange(instruction, 0, 2)
	op := getBitsRange(instruction, 10, 11)

	switch op {
	case 0: // STR
		s.write32(s.register[rb]+s.register[ro], s.register[rd])
	case 1: // STRB
		s.write8(s.register[rb]+s.register[ro], byte(s.register[rd]))
	case 2: // LDR
		s.register[rd] = s.read32(s.register[rb] + s.register[ro])
	case 3: // LDRB
		s.register[rd] = uint32(s.read8(s.register[rb] + s.register[ro]))
	}
}
