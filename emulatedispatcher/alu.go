package emulatedispatcher

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (s *CpuState) alu(instruction uint16) {
	op := getBitsRange(instruction, 6, 9)
	rs := getBitsRange(instruction, 3, 5)
	rd := getBitsRange(instruction, 0, 2)
	lhs, rhs := s.register[rd], s.register[rs]
	//lhs, rhs := s.register[rd], s.register[rs]
	switch op {
	case 0: // AND
		s.register[rd] = s.register[rd] & s.register[rs]
		s.srLogicSet(s.register[rd])
	case 1: // EOR
		s.register[rd] = s.register[rd] ^ s.register[rs]
		s.srLogicSet(s.register[rd])
	case 2: // LSL
		offset := s.register[rs] & 0xff
		if offset > 32 {
			s.sr.carry = false
			s.register[rd] = 0
		} else {
			carry := s.register[rd]&(1<<(32-offset)) > 0
			if offset > 0 {
				s.sr.carry = carry
			}
			s.register[rd] = s.register[rd] << offset
		}
		s.srLogicSet(s.register[rd])
	case 3: // LSR
		offset := s.register[rs] & 0xff
		carry := s.register[rd]&(1<<(offset-1)) > 0
		if offset > 0 {
			s.sr.carry = carry
		}
		s.register[rd] = s.register[rd] >> offset
		s.srLogicSet(s.register[rd])
	case 4: // ASR
		offset := s.register[rs] & 0xff
		val := s.register[rd]
		if offset > 32 {
			offset = 32
		}
		carry := val&(1<<(offset-1)) > 0
		if offset > 0 {
			s.sr.carry = carry
		}
		msb := (getBitsRange32(val, 31, 31) << 31)
		for i := uint(0); i < uint(offset); i++ {
			val = (val >> 1) | msb
		}
		// store the result into rd
		s.register[rd] = val
		s.srLogicSet(s.register[rd])
	case 5: // ADC
		result := uint64(lhs) + uint64(rhs) + uint64(bool2int(s.sr.carry))
		s.register[rd] = uint32(result)
		s.srArithAddSet(lhs, rhs, result)
	case 6: // SBC
		result := uint64(lhs) - uint64(rhs) - uint64(bool2int(!s.sr.carry))
		s.register[rd] = uint32(result)
		s.srArithSubSet(lhs, rhs, result)
	case 7: // ROR
		offset := s.register[rs] & 0xff
		val := s.register[rd]
		carry := (val>>(offset-1))&0b1 > 0
		if offset > 0 {
			s.sr.carry = carry
		}
		offset %= 32
		tmp0 := (val) >> (offset)
		tmp1 := (val) << (32 - (offset))
		s.register[rd] = tmp0 | tmp1
		s.srLogicSet(s.register[rd])
	case 8: // TST
		s.srLogicSet(s.register[rd] & s.register[rs])
	case 9: // NEG
		lhs = 0
		result := 0 - uint64(rhs)
		s.register[rd] = uint32(result)
		s.srArithSubSet(lhs, rhs, result)
	case 10: // CMP
		result := uint64(lhs) - uint64(rhs)
		s.srArithSubSet(lhs, rhs, result)
	case 11: // CMN
		result := uint64(lhs) + uint64(rhs)
		s.srArithAddSet(lhs, rhs, result)
	case 12: // ORR
		result := lhs | rhs
		s.register[rd] = result
		s.srLogicSet(s.register[rd])
	case 13: // MUL
		s.register[rd] = lhs * rhs
		s.sr.carry = false
		s.sr.negative = isSet(getBitsRange32(s.register[rd], 31, 31))
		s.sr.zero = s.register[rd] == 0
	case 14: // BIC
		result := lhs & ^rhs
		s.register[rd] = result
		s.srLogicSet(s.register[rd])
	case 15: // MVN
		result := ^s.register[rs]
		s.register[rd] = result
		s.srLogicSet(result)
	}
}
