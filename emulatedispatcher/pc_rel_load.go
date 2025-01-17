package emulatedispatcher

func (s *CpuState) pcRelLoad(instruction uint16) {
	rd := getBitsRange(instruction, 8, 10)
	word8 := getBitsRange(instruction, 0, 7)
	offset := word8 << 2
	wordAlignPC := (s.loc + 4) & ^uint32(3)
	s.register[rd] = s.read32(wordAlignPC + uint32(offset))
}
