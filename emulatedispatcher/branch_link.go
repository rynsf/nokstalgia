package emulatedispatcher

func (s *CpuState) branchLink(instruction uint16) {
	op := getBitsRange(instruction, 11, 11)

	switch op {
	case 0:
		offset := (int32(instruction) << 21) >> 9
		s.register[lr] = s.loc + 4
		s.register[lr] = addInt32(s.register[lr], offset)
	case 1:
		offset := getBitsRange(instruction, 0, 10)
		imm := s.register[lr] + (uint32(offset) << 1)
		s.register[lr] = (s.loc + 2)
		s.register[pc] = imm & ^uint32(1)
	}
}
