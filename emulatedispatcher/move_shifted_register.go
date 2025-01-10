package emulatedispatcher

func isSet(bin uint32) bool {
	if bin == 0 {
		return false
	}
	return true
}

func (s *CpuState) moveShiftedRegister(instruction uint16) {
	opcode := getBitsRange(instruction, 11, 12)
	offset5 := getBitsRange(instruction, 6, 10)
	rs := getBitsRange(instruction, 3, 5)
	rd := getBitsRange(instruction, 0, 2)
	switch opcode {
	case 0: // lsl
		if offset5 == 0 {
			s.register[rd] = s.register[rs]
		} else {
			// set the carry flag
			i := uint32(32 - offset5)
			carry := getBitsRange32(s.register[rs], i, i)
			if carry == 1 {
				s.sr.carry = true
			} else {
				s.sr.carry = false
			}
			s.register[rd] = s.register[rs] << offset5
		}
	case 1: // lsr
		if offset5 == 0 {
			offset5 = 32
		}
		// set the carry flag
		carry := s.register[rs]&(1<<(offset5-1)) > 0
		if offset5 > 0 {
			s.sr.carry = carry
		}
		s.register[rd] = s.register[rs] >> offset5 // logical shift right by offset
	case 2: // asr
		if offset5 == 0 {
			offset5 = 32
		}
		// set the carry flag
		val := s.register[rs]
		carry := val&(1<<(offset5-1)) > 0
		if offset5 > 0 {
			s.sr.carry = carry
		}
		msb := (getBitsRange32(val, 31, 31) << 31)
		for i := uint(0); i < uint(offset5); i++ {
			val = (val >> 1) | msb
		}
		// store the result into rd
		s.register[rd] = val
	}
	// set zero and negative flag
	s.sr.zero = s.register[rd] == 0
	s.sr.negative = isSet(getBitsRange32(s.register[rd], 31, 31))
}
