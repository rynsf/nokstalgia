package emulatedispatcher

func (s *CpuState) hiRegOpBranchEx(instruction uint16) {
	rs := getBitsRange(instruction, 3, 6)
	rd := getBitsRange(instruction, 0, 2)
	op := getBitsRange(instruction, 8, 9)

	h1 := getBitsRange(instruction, 7, 7)
	if isSet(uint32(h1)) {
		rd += 8
	}
	rsval := s.register[rs]
	rdval := s.register[rd]

	switch op {
	case 0: // ADD
		s.register[rd] = rdval + rsval
	case 1: // CMP
		result := uint64(rdval) - uint64(rsval)
		if rd != 15 {
			s.srArithSubSet(rdval, rsval, result)
		}
	case 2: // MOV
		if rs == pc {
			s.register[rd] = s.loc + 0x4
		} else {
			s.register[rd] = rsval
		}
	case 3: // BX
		rsval = rsval & ^uint32(0x1)
		s.register[pc] = rsval
	}
}
