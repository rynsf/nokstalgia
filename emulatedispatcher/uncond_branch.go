package emulatedispatcher

func (s *CpuState) uncondBranch(instruction uint16) {
	offset := (int32(instruction) << 21) >> 20
	s.register[pc] = addInt32(s.register[pc]+4, offset)
}
