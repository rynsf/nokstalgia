package emulatedispatcher

func align2(val uint32) uint32 {
	return val & ^uint32(1)
}

func (s *CpuState) pushPop(instruction uint16) {
	rlist := getBitsRange(instruction, 0, 7)
	rflag := getBitsRange(instruction, 8, 8)

	op := getBitsRange(instruction, 11, 11)

	switch op {
	case 0: // PUSH
		if rflag == 1 {
			s.register[sp] -= 4
			s.write32(s.register[sp], s.register[lr])
		}
		for i := 7; i >= 0; i-- {
			if getBitsRange(rlist, i, i) == 1 {
				s.register[sp] -= 4
				s.write32(s.register[sp], s.register[i])
			}
		}
	case 1: // POP
		for i := 0; i < 8; i++ {
			if getBitsRange(rlist, i, i) == 1 {
				s.register[i] = s.read32(s.register[sp])
				s.register[sp] += 4
			}
		}
		if rflag == 1 {
			s.register[pc] = s.read32(s.register[sp])
			s.register[pc] = align2(s.register[pc])
			s.register[sp] += 4
		}
	}
}
