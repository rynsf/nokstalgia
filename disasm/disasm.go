package main

import (
	"fmt"
	"io"
	"os"
)

var out io.Writer = os.Stdout

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
*/
func disassembler(pc int, codebuffer []byte) {
	instruction := bytesToHalfword(codebuffer[pc], codebuffer[pc+1])
	first3Bit := getBitsRange(instruction, 13, 15)
	fmt.Printf("0x%04X ", pc)
	fmt.Printf("%04X ", instruction)
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
		mnemonics := []string{"STR", "LDR", "STRB", "LDRB"}
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
			fmt.Fprint(out, "load address")
		case 1: // add offset to stack pointer & push/pop register
			bit10 := getBitsRange(instruction, 10, 10)
			switch bit10 {
			case 0: // add offest to stack pointer
				fmt.Fprint(out, "add offest to stack pointer")
			case 1: // push/pop register
				fmt.Fprint(out, "push/pop register")
			}
		}

	case 6: // multiple load/store & conditional branch & software interrupt
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // multiple load/store
			fmt.Fprint(out, "multiple load/store")
		case 1:
			bitCond := getBitsRange(instruction, 8, 11)
			switch bitCond {
			case 15: // software interrupt
				fmt.Fprint(out, "software interrupt")
			default: // conditional branch
				fmt.Fprint(out, "conditional branch")
			}
		}

	case 7: // unconditional branch / long branch with link
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // unconditional branch
			fmt.Fprint(out, "unconditional branch")
		case 1: // long branch with link
			fmt.Fprint(out, "long branch with link")
		}
	}
	fmt.Println()
}

func main() {
	bin, err := os.ReadFile("./outputcode.bin")
	if err != nil {
		panic(err)
	}
	for pc := 0; pc < len(bin); pc += 2 {
		disassembler(pc, bin)
	}
}
