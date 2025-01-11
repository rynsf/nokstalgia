package emulatedispatcher

func (s *CpuState) movCmpAddSubIm(instruction uint16) {
	opcode := getBitsRange(instruction, 11, 12)
	rd := getBitsRange(instruction, 8, 10)
	offset8 := getBitsRange(instruction, 0, 7)
	lhs := s.register[rd]
	rhs := uint32(offset8)
	switch opcode {
	case 0: // MOV Rd, #Offset8
		s.register[rd] = uint32(rhs)
		s.srLogicSet(s.register[rd])
	case 1: // CMP Rd, #Offset8
		result := uint64(lhs) - uint64(rhs)
		s.srArithSubSet(lhs, rhs, result)
	case 2: // ADD Rd, #Offset8
		result := uint64(lhs) + uint64(rhs)
		s.register[rd] = uint32(result)
		s.srArithAddSet(lhs, rhs, result)
	case 3: // SUB Rd, #Offset8
		result := uint64(lhs) - uint64(rhs)
		s.register[rd] = uint32(result)
		s.srArithSubSet(lhs, rhs, result)
	}
}
