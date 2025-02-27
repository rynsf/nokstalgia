package emulatedispatcher

import (
	"encoding/binary"
	"fmt"
)

type cpsr struct {
	overflow bool
	carry    bool
	zero     bool
	negative bool
}

type CpuState struct {
	register     [16]uint32
	sr           cpsr
	loc          uint32
	ram          []byte
	ramBaseAdr   uint32
	dynamicRam   []byte
	dynRamBase   uint32
	dynRamLen    uint32
	ramLen       uint32
	flash        []byte
	flashBaseAdr uint32
	flashLen     uint32
}

const (
	sp = 13
	lr = 14
	pc = 15
)

func (s CpuState) read32(address uint32) uint32 {
	address = address & ^uint32(3) // word aligned
	var result uint32
	result = result | (uint32(s.read(address)) << 24)
	result = result | (uint32(s.read(address+1)) << 16)
	result = result | (uint32(s.read(address+2)) << 8)
	result = result | uint32(s.read(address+3))
	return result
}

func (s CpuState) read16(address uint32) uint32 {
	var result uint32
	result = result | (uint32(s.read(address)) << 8)
	result = result | uint32(s.read(address+1))
	return result
}

func (s CpuState) read8(address uint32) byte {
	return s.read(address)
}

func (s CpuState) read(address uint32) byte {
	if Debug {
		fmt.Printf("read at: %X from: %X\n", address, s.loc)
	}
	if address >= s.ramBaseAdr && address < s.ramBaseAdr+s.ramLen {
		absoluteAddress := address - s.ramBaseAdr
		return s.ram[absoluteAddress]
	}
	if address >= s.flashBaseAdr && address < s.flashBaseAdr+s.flashLen {
		absoluteAddress := address - s.flashBaseAdr
		return s.flash[absoluteAddress]
	}
	if address >= s.dynRamBase && address < s.dynRamBase+s.dynRamLen {
		absoluteAddress := address - s.dynRamBase
		return s.dynamicRam[absoluteAddress]
	}
	if Debug {
		fmt.Printf("Out of bound read at: %X from: %X\n", address, s.loc)
	}
	return 0x0
}

func (s *CpuState) write32(address, data uint32) {
	address = address & ^uint32(3) // word aligned
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, data)
	for i := uint32(0); i < 4; i++ {
		s.write(address+i, bytes[i])
	}
}

func (s *CpuState) write16(address uint32, data uint16) {
	byte1 := (data >> 8) & uint16(0xff)
	byte2 := data & uint16(0xff)
	s.write(address, byte(byte1))
	s.write(address+1, byte(byte2))
}

func (s *CpuState) write8(address uint32, data byte) {
	s.write(address, data)
}

func (s *CpuState) write(address uint32, data byte) {
	if Debug {
		fmt.Printf("write at: %X from: %X\n", address, s.loc)
	}
	if address >= s.ramBaseAdr && address < s.ramBaseAdr+s.ramLen {
		absoluteAddress := address - s.ramBaseAdr
		s.ram[absoluteAddress] = data
		return
	}
	if address >= s.flashBaseAdr && address < s.flashBaseAdr+s.flashLen {
		absoluteAddress := address - s.flashBaseAdr
		s.flash[absoluteAddress] = data
		return
	}
	if address >= s.dynRamBase && address < s.dynRamBase+s.dynRamLen {
		absoluteAddress := address - s.dynRamBase
		s.dynamicRam[absoluteAddress] = data
		return
	}
	if Debug {
		fmt.Printf("Out of bound write at: %X from: %X\n", address, s.loc)
	}
}
