package emulatedispatcher

func (s *CpuState) srArithAddSet(lhs, rhs uint32, result uint64) {
	s.sr.zero = uint32(result) == 0
	s.sr.negative = isSet(getBitsRange32(uint32(result), 31, 31))
	s.sr.carry = isAddCarry(result)
	s.sr.overflow = isAddOverflow(lhs, rhs, uint32(result))
}

func (s *CpuState) srArithSubSet(lhs, rhs uint32, result uint64) {
	s.sr.zero = uint32(result) == 0
	s.sr.negative = isSet(getBitsRange32(uint32(result), 31, 31))
	s.sr.carry = isSubCarry(result)
	s.sr.overflow = isSubOverflow(lhs, rhs, uint32(result))
}

func (s *CpuState) addSub(instruction uint16) {
	op := getBitsRange(instruction, 9, 10)
	rn_offset3 := getBitsRange(instruction, 6, 8)
	rs := getBitsRange(instruction, 3, 5)
	rd := getBitsRange(instruction, 0, 2)
	lhs := s.register[rs]
	rhs := s.register[rn_offset3]
	switch op {
	case 0: // ADD Rd, Rs, Rn
		result := uint64(lhs) + uint64(rhs)
		s.register[rd] = uint32(result)
		s.srArithAddSet(lhs, rhs, result)
	case 1: // SUB Rd, Rs, Rn
		result := uint64(lhs) - uint64(rhs)
		s.register[rd] = uint32(lhs - rhs)
		s.srArithSubSet(lhs, rhs, result)
	case 2: // ADD Rd, Rs, #Offset3
		rhs = uint32(rn_offset3)
		result := uint64(lhs) + uint64(rhs)
		s.register[rd] = uint32(result)
		s.srArithAddSet(lhs, rhs, result)
	case 3: // SUB Rd, Rs, #Offset3
		rhs = uint32(rn_offset3)
		result := uint64(lhs) - uint64(rhs)
		s.register[rd] = uint32(result)
		s.srArithSubSet(lhs, rhs, result)
	}
}
