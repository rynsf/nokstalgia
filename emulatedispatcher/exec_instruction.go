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

func (s *CpuState) execInstruction(instruction uint16) {
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
		opbits := getBitsRange(instruction, 11, 12)
		switch opbits {
		case 0:
			bit10 := getBitsRange(instruction, 10, 10)
			switch bit10 {
			case 0:
				s.alu(instruction)
			case 1:
				s.hiRegOpBranchEx(instruction)
			}
		case 1:
			s.pcRelLoad(instruction)
		case 2, 3:
			bit9 := getBitsRange(instruction, 9, 9)
			switch bit9 {
			case 0: // load/store with register offset
				s.loadStoreRegOffset(instruction)
			case 1: // load/store sign-extended byte/halfword
				s.loadStoreSBH(instruction)
			}
		}
	case 3: // load/store with immediate offset
		s.loadStoreImmOffset(instruction)
	case 4:
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // load/store halfword
			s.loadStoreHalfword(instruction)
		case 1: // SP-relative load/store
			s.spRelLoadStore(instruction)
		}
	case 5:
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // load address
			s.loadAddress(instruction)
		case 1:
			bit10 := getBitsRange(instruction, 10, 10)
			switch bit10 {
			case 0: // add offset to stack pointer
				s.addOffsetSP(instruction)
			case 1: // push/pop instruction
				s.pushPop(instruction)
			}
		}
	case 6:
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // multiple load/store
			s.multipleLoadStore(instruction)
		case 1:
			op := getBitsRange(instruction, 8, 11)
			switch op {
			case 15: // SWI
			default: // conditional branch
				s.conditionalBranch(instruction)
			}
		}
	case 7:
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // unconditional branch
			s.uncondBranch(instruction)
		case 1: // long branch with link
			s.branchLink(instruction)
		}
	}
}
