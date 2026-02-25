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

func isArmALU(instruction uint32) bool {
	return instruction&0b0000_1100_0000_0000_0000_0000_0000_0000 == 0
}

func (s *CpuState) execArmALU(instruction uint32) {
}

func isArmBx(instruction uint32) bool {
	return instruction&0b0000_1111_1111_1111_1111_1111_1111_0000 == 0b0000_0001_0010_1111_1111_1111_0001_0000
}

func (s *CpuState) execArmBx(instruction uint32) {
	rn := getBitsRange32(instruction, 0, 3)
	rnVal := s.register[rn]

	if rnVal&0x1 == 0 {
		s.thumb = false
	} else {
		s.thumb = true
	}

	rnVal = rnVal & ^uint32(0x1)
	s.register[pc] = rnVal
}

func isArmMull(instruction uint32) bool {
	return instruction&0b0000_1111_1000_0000_0000_0000_1111_0000 == 0b0000_0000_1000_0000_0000_0000_1001_0000
}

func (s *CpuState) execArmMull(instruction uint32) {
	opcode := getBitsRange32(instruction, 20, 22)
	switch opcode {
	case 0: // UMALL
		rdhi := getBitsRange32(instruction, 16, 19)
		rdlo := getBitsRange32(instruction, 12, 15)
		rs := getBitsRange32(instruction, 8, 11)
		rm := getBitsRange32(instruction, 0, 3)

		result := uint64(s.register[rs]) * uint64(s.register[rm])
		s.register[rdhi], s.register[rdlo] = uint32(result>>32), uint32(result)
		if getBitsRange32(instruction, 20, 20) == 1 {
			s.sr.zero = result == 0
			s.sr.negative = isSet(getBitsRange32(s.register[rdhi], 31, 31))
		}
	}
}

func (s *CpuState) execArmInstruction(instruction uint32) {
	// docode and execute arm instruction
	// only implementing the instructions that the phones actually executes
	// ARM mode is only used in FD_MUL function for multiplication of
	// floating point numbers, which is in turned used by sqrt in link5

	switch {
	case isArmBx(instruction): // bx
		s.execArmBx(instruction)
	case isArmMull(instruction): // mull
		s.execArmMull(instruction)
	case isArmALU(instruction): // adr, adds and mov
		s.execArmALU(instruction)
	}

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
