package emulatedispatcher

func (s *CpuState) multipleLoadStore(instruction uint16) {
	rlist := getBitsRange(instruction, 0, 7)
	rb := getBitsRange(instruction, 8, 10)
	op := getBitsRange(instruction, 11, 11)

	switch op {
	case 0: // STMIA
		for i := 0; i < 8; i++ {
			if getBitsRange(rlist, i, i) == 1 {
				s.write32(s.register[rb], s.register[i])
				s.register[rb] += 4
			}
		}
	case 1: // LDMIA
		basePtr := s.register[rb]
		rbInRlist := false
		for i := 0; i < 8; i++ {
			if getBitsRange(rlist, i, i) == 1 {
				s.register[i] = s.read32(basePtr)
				basePtr += 4
				if int(rb) == i {
					rbInRlist = true
				}
			}
		}
		if !rbInRlist {
			s.register[rb] = basePtr
		}
	}
}
