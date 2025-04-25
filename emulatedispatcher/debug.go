package emulatedispatcher

import (
	"bufio"
	"fmt"
	"os"
)

var Debug bool // when set, step through instructions and print debug info

func (s *CpuState) GetRam() []byte {
	return s.ram
}

func (s *CpuState) GetDynamicRam() []byte {
	return s.dynamicRam
}

func (s *CpuState) GetReg(i int) uint32 {
	return s.register[i]
}

func (s *CpuState) SetReg(val uint32, i int) {
	s.register[i] = val
}

func (s *CpuState) ResetReg() {
	for i := 0; i < 16; i++ {
		s.register[i] = 0x0
	}
}

func (s *CpuState) FillMem(data []byte, base int) int {
	for i := range data {
		s.ram[base+i] = data[i]
	}
	return base
}

func (s *CpuState) DumbState() {
	fmt.Println("Registers: ")
	for i := 0; i < 13; i++ {
		fmt.Printf("%d: %X\t", i, s.register[i])
	}
	fmt.Println()
	fmt.Printf("LR: %X\n", s.register[lr])
	fmt.Printf("SP: %X\n", s.register[sp])
	fmt.Printf("PC: %X\n", s.register[pc])
	fmt.Printf("v: %v\t", s.sr.overflow)
	fmt.Printf("c: %v\t", s.sr.carry)
	fmt.Printf("z: %v\t", s.sr.zero)
	fmt.Printf("n: %v\t", s.sr.negative)
	fmt.Println()
}

func (s *CpuState) step() {
	fmt.Print("Press enter to step...")
	reader := bufio.NewReader(os.Stdin)
	i, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if i == "c\n" {
		Debug = false
	}
	s.DumbState()
}
