package emulatedispatcher

func addInt32(u uint32, i int32) uint32 {
	if i > 0 {
		u += uint32(i)
	} else {
		u -= uint32(-i)
	}
	return u
}

func (s *CpuState) check(cond uint16) bool {
	switch cond {
	case 0: // ==
		if s.sr.zero {
			return true
		}
	case 1: // !=
		if !s.sr.zero {
			return true
		}
	case 2: // carry
		if s.sr.carry {
			return true
		}
	case 3: // not carry
		if !s.sr.carry {
			return true
		}
	case 4: // minus
		if s.sr.negative {
			return true
		}
	case 5: // plus
		if !s.sr.negative {
			return true
		}
	case 6: // overflow
		if s.sr.overflow {
			return true
		}
	case 7: // not-overflow
		if !s.sr.overflow {
			return true
		}
	case 8: // higher
		if s.sr.carry && !s.sr.zero {
			return true
		}
	case 9: // not carry
		if !s.sr.carry && s.sr.zero {
			return true
		}
	case 10: // >=
		if s.sr.negative == s.sr.overflow {
			return true
		}
	case 11: // <
		if s.sr.negative != s.sr.overflow {
			return true
		}
	case 12: // >
		if !s.sr.zero && (s.sr.negative == s.sr.overflow) {
			return true
		}
	case 13: // <=
		if s.sr.zero && (s.sr.negative != s.sr.overflow) {
			return true
		}
	}
	return false
}

func (s *CpuState) conditionalBranch(instruction uint16) {
	cond := getBitsRange(instruction, 8, 11)

	if s.check(cond) {
		offset := int32(int8(byte(instruction & 0b1111_1111)))
		s.register[pc] = addInt32(s.loc+4, offset*2)
	}
}
