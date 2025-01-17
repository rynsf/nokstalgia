package emulatedispatcher

import (
	"fmt"
)

func (s *CpuState) DumbState() {
	fmt.Println("Registers: ")
	for i := 0; i < 13; i++ {
		fmt.Printf("%d: %d\t", i, s.register[i])
	}
	fmt.Println()
	fmt.Println("LR: ", s.register[lr])
	fmt.Println("SP: ", s.register[sp])
	fmt.Println("PC: ", s.register[pc])
	fmt.Printf("v: %v\t", s.sr.overflow)
	fmt.Printf("c: %v\t", s.sr.carry)
	fmt.Printf("z: %v\t", s.sr.zero)
	fmt.Printf("n: %v\t", s.sr.negative)
	fmt.Println()
}

func (s *CpuState) subroutineIsRunning() bool {
	return s.register[pc] != 0
}

func InitCpu(r [16]uint32, v, c, z, n bool, ram []byte, ramBase, ramLen uint32, flash []byte, flashBase, flashLen uint32) CpuState {
	return CpuState{
		register:     r,
		sr:           cpsr{v, c, z, n},
		loc:          0,
		ram:          ram,
		ramBaseAdr:   ramBase,
		ramLen:       ramLen,
		flash:        flash,
		flashBaseAdr: flashBase,
		flashLen:     flashLen,
	}
}

func (s *CpuState) RunSubroutine() {
	for s.subroutineIsRunning() {
		instruction := uint16(s.read16(s.register[pc]))
		s.loc = s.register[pc]
		s.register[pc] += 2
		s.execInstruction(instruction)
	}
}
