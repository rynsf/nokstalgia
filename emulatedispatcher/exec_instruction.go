package emulatedispatcher

func bytesToHalfword(byteHi, byteLo byte) uint16 {
	return (uint16(byteHi) << 8) + uint16(byteLo)
}

func getBitsRange32(w uint32, start, end uint32) uint32 {
	if start > end || start >= 32 || end >= 32 {
		panic("invalid range")
	}
	rightChoped := w << (31 - end)
	leftChoped := rightChoped >> (31 - (end - start))
	return leftChoped
}

func getBitsRange(w uint16, start, end int) uint16 {
	if start > end || start >= 16 || end >= 16 {
		panic("invalid range")
	}
	rightChoped := w << (15 - end)
	leftChoped := rightChoped >> (15 - (end - start))
	return leftChoped
}

func (s *CpuState) ExecInstruction(instruction uint16) {
	// decode the instrution
	// call the appropriate function that implements the instruction

	//decode instruction
	first3Bit := getBitsRange(instruction, 13, 15)
	switch first3Bit {
	case 0: // move shifted register and add/substract
		opcode := getBitsRange(instruction, 11, 12)
		switch opcode {
		case 3: // add / substract
			s.addSub(instruction)
		default: // move shifted register
			s.moveShiftedRegister(instruction)
		}
	case 1:
		s.movCmpAddSubIm(instruction)
	case 2:
	case 3:
	case 4:
	case 5:
	case 6:
	case 7:
	}
}
