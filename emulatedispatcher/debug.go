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

func (s *CpuState) SetReg(val uint32, i int) {
	s.register[i] = val
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
	_, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	s.DumbState()
}
