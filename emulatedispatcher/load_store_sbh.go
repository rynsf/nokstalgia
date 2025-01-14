package emulatedispatcher

func (s *CpuState) loadStoreSBH(instruction uint16) {
	ro := getBitsRange(instruction, 6, 8)
	rb := getBitsRange(instruction, 3, 5)
	rd := getBitsRange(instruction, 0, 2)

	op := getBitsRange(instruction, 10, 11)

	switch op {
	case 0: // STRH
		s.write16(s.register[rb]+s.register[ro], uint16(s.register[rd]))
	case 1: // LDSB
		value := int32(s.read8(s.register[rb] + s.register[ro]))
		value = (value << 24) >> 24
		s.register[rd] = uint32(value)
	case 2: // LDRH
		s.register[rd] = uint32(s.read16(s.register[rb] + s.register[ro]))
	case 3: // LDSH
		addr := s.register[rb] + s.register[ro]
		val := s.read16(addr)
		if addr%2 == 1 {
			val = ((val & 0xff) << 24) | ((val & 0xff) << 16) | ((val & 0xff) << 8) | val
		}
		value := int32(val)
		value = (value << 16) >> 16
		s.register[rd] = uint32(value)
	}
}
