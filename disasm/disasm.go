package main

import (
	"fmt"
	"io"
	"os"
)

var out io.Writer = os.Stdout

// twosComplementToInt takes in bytes and the length of those bytes
// and returns signed integer represented by the input.
func twosComplementToInt(bytes uint32, bitLen int) int {
	if bytes&(1<<(bitLen-1)) != 0 {
		return int(bytes) - (1 << bitLen)
	} else {
		return int(bytes)
	}
}

/*
function: bytesToHalfword

parameter:

	byteHi: high byte
	byteLo: low byte

return: half word formed by concatnating the two bytes
*/
func bytesToHalfword(byteHi, byteLo byte) uint16 {
	return (uint16(byteHi) << 8) + uint16(byteLo)
}

/*
function: getBitsRange

parameter:

	w: the half word from with bytes have to be taken.
	start: start index of the bit range to return, inclusive, index starts from 0
	end: end index of the bit range to return, inclusive, index starts form 0

return: the byte range specified by the range
*/
func getBitsRange(w uint16, start, end int) uint16 {
	if start > end || start >= 16 || end >= 16 {
		panic("invalid range")
	}
	rightChoped := w << (15 - end)
	leftChoped := rightChoped >> (15 - (end - start))
	return leftChoped
}

/*
function: disassembler
Prints the assembly for instruction

parameter:

	pc: index that points to the code to disassemble
	codebuffer: byte slice of binary code to disassemble

returns: the number of bytes read
*/
func disassembler(pc int, codebuffer []byte) int {
	instruction := bytesToHalfword(codebuffer[pc], codebuffer[pc+1])
	first3Bit := getBitsRange(instruction, 13, 15)
	fmt.Fprintf(out, "0x%04X ", pc)
	fmt.Fprintf(out, "0x%04X ", instruction)
	switch first3Bit {
	case 0: // move shifted register and add/substract
		opBytes := getBitsRange(instruction, 11, 12)
		switch opBytes {
		case 0: // lsl
			offset5 := getBitsRange(instruction, 6, 10)
			rs := getBitsRange(instruction, 3, 5)
			rd := getBitsRange(instruction, 0, 2)
			fmt.Fprintf(out, "LSL R%d, R%d, #%d", rd, rs, offset5)
		case 1: // lsr
			offset5 := getBitsRange(instruction, 6, 10)
			rs := getBitsRange(instruction, 3, 5)
			rd := getBitsRange(instruction, 0, 2)
			fmt.Fprintf(out, "LSR R%d, R%d, #%d", rd, rs, offset5)
		case 2: //asr
			offset5 := getBitsRange(instruction, 6, 10)
			rs := getBitsRange(instruction, 3, 5)
			rd := getBitsRange(instruction, 0, 2)
			fmt.Fprintf(out, "ASR R%d, R%d, #%d", rd, rs, offset5)
		case 3: // add/sub
			immediateFlag := getBitsRange(instruction, 10, 10)
			opcode := getBitsRange(instruction, 9, 9)
			rn_offset3 := getBitsRange(instruction, 6, 8)
			rs := getBitsRange(instruction, 3, 5)
			rd := getBitsRange(instruction, 0, 2)
			switch opcode {
			case 0: // add
				switch immediateFlag {
				case 0: //register
					fmt.Fprintf(out, "ADD R%d, R%d, R%d", rd, rs, rn_offset3)
				case 1: // immediate
					fmt.Fprintf(out, "ADD R%d, R%d, #%d", rd, rs, rn_offset3)
				}
			case 1: //sub
				switch immediateFlag {
				case 0: //register
					fmt.Fprintf(out, "SUB R%d, R%d, R%d", rd, rs, rn_offset3)
				case 1: // immediate
					fmt.Fprintf(out, "SUB R%d, R%d, #%d", rd, rs, rn_offset3)
				}
			}
		}

	case 1: // move/compare/add/sub immediate
		opBytes := getBitsRange(instruction, 11, 12)
		rd := getBitsRange(instruction, 8, 10)
		offset8 := getBitsRange(instruction, 0, 7)
		switch opBytes {
		case 0: // mov
			fmt.Fprintf(out, "MOV R%d #%d", rd, offset8)
		case 1: // cmp
			fmt.Fprintf(out, "CMP R%d #%d", rd, offset8)
		case 2: // add
			fmt.Fprintf(out, "ADD R%d #%d", rd, offset8)
		case 3: // sub
			fmt.Fprintf(out, "SUB R%d #%d", rd, offset8)
		}

	case 2: // alu operations / high register operation / branch exchange / pc relative load / load / store with register offset / load store sign extended byte/halfword
		opBytes := getBitsRange(instruction, 11, 12)
		switch opBytes {
		case 0:
			bit10 := getBitsRange(instruction, 10, 10)
			switch bit10 {
			case 0: // alu operation
				opcode := getBitsRange(instruction, 6, 9)
				rs := getBitsRange(instruction, 3, 5)
				rd := getBitsRange(instruction, 0, 2)
				mnemonics := []string{"AND", "EOR", "LSL", "LSR", "ASR", "ADC", "SBC", "ROR", "TST", "NEG", "CMP", "CMN", "ORR", "MUL", "BIC", "MVN"}
				fmt.Fprintf(out, "%s R%d R%d", mnemonics[opcode], rd, rs)
			case 1: // high register operation/ branch exchange
				opcode := getBitsRange(instruction, 8, 9)
				h1 := getBitsRange(instruction, 7, 7)
				h2 := getBitsRange(instruction, 6, 6)
				rs := getBitsRange(instruction, 3, 5)
				rd := getBitsRange(instruction, 0, 2)
				h_r := []string{"R", "H"}
				switch opcode {
				case 0:
					fmt.Fprintf(out, "ADD %s%d, %s%d", h_r[h1], rd, h_r[h2], rs)
				case 1:
					fmt.Fprintf(out, "CMP %s%d, %s%d", h_r[h1], rd, h_r[h2], rs)
				case 2:
					fmt.Fprintf(out, "MOV %s%d, %s%d", h_r[h1], rd, h_r[h2], rs)
				case 3:
					fmt.Fprintf(out, "BX %s%d", h_r[h2], rs)
				}
			}
		case 1: // pc relative laod
			rd := getBitsRange(instruction, 8, 10)
			word8 := getBitsRange(instruction, 0, 7)
			fmt.Fprintf(out, "LDR R%d, [PC, #%d]", rd, word8)
		case 2, 3: // load/store with regiter offset & load/store sign-extended byte/halfword
			bit9 := getBitsRange(instruction, 9, 9)
			switch bit9 {
			case 0: // load/store with regiter offset
				lflag := getBitsRange(instruction, 11, 11)
				bflag := getBitsRange(instruction, 10, 10)
				ro := getBitsRange(instruction, 6, 8)
				rb := getBitsRange(instruction, 3, 5)
				rd := getBitsRange(instruction, 0, 2)
				opcode := (lflag * 2) + bflag
				mnemonics := []string{"STR", "STRB", "LDR", "LDRB"}
				fmt.Fprintf(out, "%s R%d [R%d, R%d]", mnemonics[opcode], rd, rb, ro)
			case 1: // load/store sign-extended byte/halfword
				hflag := getBitsRange(instruction, 11, 11)
				sflag := getBitsRange(instruction, 10, 10)
				ro := getBitsRange(instruction, 6, 8)
				rb := getBitsRange(instruction, 3, 5)
				rd := getBitsRange(instruction, 0, 2)
				opcode := (sflag * 2) + hflag
				mnemonics := []string{"STRH", "LDRH", "LDSB", "LDSH"}
				fmt.Fprintf(out, "%s R%d [R%d, R%d]", mnemonics[opcode], rd, rb, ro)
			}
		}

	case 3: // load/store with immediate offset
		bflag := getBitsRange(instruction, 12, 12)
		lflag := getBitsRange(instruction, 11, 11)
		offset5 := getBitsRange(instruction, 6, 10)
		offset := offset5
		if bflag == 0 {
			offset = offset5 << 2
		}
		rb := getBitsRange(instruction, 3, 5)
		rd := getBitsRange(instruction, 0, 2)
		mnemonics := []string{"STR", "STRB", "LDR", "LDRB"}
		opcode := (lflag * 2) + bflag
		fmt.Fprintf(out, "%s R%d, [R%d, #%d]", mnemonics[opcode], rd, rb, offset)

	case 4: // load/store half word & sp relative load store
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // load/store half word
			lflag := getBitsRange(instruction, 11, 11)
			offset5 := getBitsRange(instruction, 6, 10)
			offset := offset5 << 1
			rb := getBitsRange(instruction, 3, 5)
			rd := getBitsRange(instruction, 0, 2)
			mnemonic := []string{"STRH", "LDRH"}[lflag]
			fmt.Fprintf(out, "%s R%d, [R%d, #%d]", mnemonic, rd, rb, offset)
		case 1: // sp relative load/store
			lflag := getBitsRange(instruction, 11, 11)
			rd := getBitsRange(instruction, 8, 10)
			word8 := getBitsRange(instruction, 0, 7)
			offset := word8 << 2
			mnemonic := []string{"STR", "LDR"}[lflag]
			fmt.Fprintf(out, "%s R%d, [SP, #%d]", mnemonic, rd, offset)
		}

	case 5:
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // load address
			pc_sp := getBitsRange(instruction, 11, 11)
			base := []string{"PC", "SP"}[pc_sp]
			rd := getBitsRange(instruction, 8, 10)
			word8 := getBitsRange(instruction, 0, 7)
			offset := word8 << 2
			fmt.Fprintf(out, "ADD R%d, %s, #%d", rd, base, offset)
		case 1: // add offset to stack pointer & push/pop register
			bit10 := getBitsRange(instruction, 10, 10)
			switch bit10 {
			case 0: // add offest to stack pointer
				signFlag := getBitsRange(instruction, 7, 7)
				sign := []string{"", "-"}[signFlag]
				sword7 := getBitsRange(instruction, 0, 6)
				offset := sword7 << 2
				fmt.Fprintf(out, "ADD SP, #%s%d", sign, offset)
			case 1: // push/pop register
				lflag := getBitsRange(instruction, 11, 11)
				rflag := getBitsRange(instruction, 8, 8)
				rlist := getBitsRange(instruction, 0, 7)
				opcode := (lflag * 2) + rflag
				rlistStr := ""
				sep := ""
				for bitIndex := 0; bitIndex < 8; bitIndex++ {
					if getBitsRange(rlist, bitIndex, bitIndex) == 1 {
						rlistStr += fmt.Sprintf("%sR%d", sep, bitIndex)
						sep = ", "
					}
				}
				switch opcode {
				case 0:
					fmt.Fprintf(out, "PUSH { %s }", rlistStr)
				case 1:
					fmt.Fprintf(out, "PUSH { %s, LR }", rlistStr)
				case 2:
					fmt.Fprintf(out, "POP { %s }", rlistStr)
				case 3:
					fmt.Fprintf(out, "POP { %s, PC }", rlistStr)
				}
			}
		}

	case 6: // multiple load/store & conditional branch & software interrupt
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // multiple load/store
			lflag := getBitsRange(instruction, 11, 11)
			rb := getBitsRange(instruction, 8, 10)
			rlist := getBitsRange(instruction, 0, 7)
			mnemonic := []string{"STMIA", "LDMIA"}[lflag]
			rlistStr := ""
			sep := ""
			for bitIndex := 0; bitIndex < 8; bitIndex++ {
				if getBitsRange(rlist, bitIndex, bitIndex) == 1 {
					rlistStr += fmt.Sprintf("%sR%d", sep, bitIndex)
					sep = ", "
				}
			}
			fmt.Fprintf(out, "%s R%d!, { %s }", mnemonic, rb, rlistStr)
		case 1:
			bitCond := getBitsRange(instruction, 8, 11)
			switch bitCond {
			case 15: // software interrupt
				value8 := getBitsRange(instruction, 0, 7)
				fmt.Fprintf(out, "SWI %d", value8)
			default: // conditional branch
				cond := getBitsRange(instruction, 8, 11)
				soffset8 := getBitsRange(instruction, 0, 7)
				offset := twosComplementToInt(uint32(soffset8<<1), 9)
				mnemonic := []string{"BEQ", "BNE", "BCS", "BCC", "BMI", "BPL", "BVS", "BVC", "BHI", "BLS", "BGE", "BLT", "BGT", "BLE"}[cond]
				label := fmt.Sprintf("LAB_%X", pc+offset)
				fmt.Fprintf(out, "%s %s", mnemonic, label)
			}
		}

	case 7: // unconditional branch / long branch with link
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // unconditional branch
			offset11 := getBitsRange(instruction, 0, 10)
			offset := twosComplementToInt(uint32(offset11), 11) << 1
			label := fmt.Sprintf("LAB_%X", pc+offset)
			fmt.Fprintf(out, "B %s", label)
		case 1: // long branch with link
			offsetHigh := getBitsRange(instruction, 0, 10)
			nextInstruction := bytesToHalfword(codebuffer[pc+2], codebuffer[pc+3])
			offsetLow := getBitsRange(nextInstruction, 0, 10)
			offsetH := twosComplementToInt(uint32(offsetHigh), 11)
			offsetL := twosComplementToInt(uint32(offsetLow), 11)
			offset := (offsetH << 12) + (offsetL << 1)
			label := fmt.Sprintf("LAB_%X", pc+offset)
			fmt.Fprintf(out, "0x%04X ", nextInstruction)
			fmt.Fprintf(out, "BL %s\n", label)
			return 4
		}
	}
	fmt.Println()
	return 2
}

func main() {
	bin, err := os.ReadFile("./outputcode.bin")
	if err != nil {
		panic(err)
	}
	pc := 0
	for pc < len(bin) {
		pc += disassembler(pc, bin)
	}
}
