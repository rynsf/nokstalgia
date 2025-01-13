package emulatedispatcher

type cpsr struct {
	overflow bool
	carry    bool
	zero     bool
	negative bool
}

type CpuState struct {
	register     [16]uint32
	sr           cpsr
	ram          []byte
	ramBaseAdr   uint32
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

func (s CpuState) read(address uint32) byte {
	if address >= s.ramBaseAdr && address < s.ramBaseAdr+s.ramLen {
		absoluteAddress := address - s.ramBaseAdr
		return s.ram[absoluteAddress]
	}
	if address >= s.flashBaseAdr && address < s.flashBaseAdr+s.flashLen {
		absoluteAddress := address - s.flashBaseAdr
		return s.flash[absoluteAddress]
	}
	return 0x0
}

func (s *CpuState) write(address uint32, data byte) {
	if address >= s.ramBaseAdr && address < s.ramBaseAdr+s.ramLen {
		absoluteAddress := address - s.ramBaseAdr
		s.ram[absoluteAddress] = data
	}
	if address >= s.flashBaseAdr && address < s.flashBaseAdr+s.flashLen {
		absoluteAddress := address - s.flashBaseAdr
		s.flash[absoluteAddress] = data
	}
}
