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
			fmt.Fprint(out, "lsr")
		case 2: //asr
			fmt.Fprint(out, "asr")
		case 3: // add/sub
			fmt.Fprint(out, "add/sub")
		}

	case 1: // move/compare/add/sub immediate
		opBytes := getBitsRange(instruction, 11, 12)
		switch opBytes {
		case 0:
			fmt.Fprint(out, "mov")
		case 1:
			fmt.Fprint(out, "cmp")
		case 2:
			fmt.Fprint(out, "add")
		case 3:
			fmt.Fprint(out, "sub")
		}

	case 2: // alu operations / high register operation / branch exchange / pc relative load / load / store with register offset / load store sign extended byte/halfword
		opBytes := getBitsRange(instruction, 11, 12)
		switch opBytes {
		case 0:
			bit10 := getBitsRange(instruction, 10, 10)
			switch bit10 {
			case 0: // alu operation
				fmt.Fprint(out, "alu op")
			case 1: // high register operation/ branch exchange
				fmt.Fprint(out, "hi reg")
			}
		case 1: // pc relative laod
			fmt.Fprint(out, "PC relative load")
		case 2, 3: // load/store with regiter offset & load/store sign-extended byte/halfword
			bit9 := getBitsRange(instruction, 9, 9)
			switch bit9 {
			case 0: // load/store with regiter offset
				fmt.Fprint(out, "load/store with regiter offset")
			case 1: // load/store sign-extended byte/halfword
				fmt.Fprint(out, "load/store sign-extended byte/halfword")
			}
		}

	case 3: // load/store with immediate offset
		fmt.Fprint(out, "load/store with immediate offset")

	case 4: // load/store half word & sp relative load store
		bit12 := getBitsRange(instruction, 12, 12)
		switch bit12 {
		case 0: // load/store half word
			fmt.Fprint(out, "load/store half word")
		case 1: // sp relative load/store
			fmt.Fprint(out, "sp relative load/store")
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
