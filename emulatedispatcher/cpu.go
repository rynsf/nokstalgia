package emulatedispatcher

type cpsr struct {
	overflow bool
	carry    bool
	zero     bool
	negative bool
}

type CpuState struct {
	register [16]uint32
	cr       cpsr
	ram      []byte
	flash    []byte
}

const (
	sp = 13
	lr = 14
	pc = 15
)

func (s CpuState) read(address uint32) byte {
	if address >= 0x100000 && address <= 0x10FFFF {
		baseAddress := uint32(0x100000)
		absoluteAddress := address - baseAddress
		return s.ram[absoluteAddress]
	}
	if address >= 0x200000 && address <= 0x2fffff {
		baseAddress := uint32(0x200000)
		absoluteAddress := address - baseAddress
		return s.flash[absoluteAddress]
	}
	return 0x0
}

func (s *CpuState) write(address uint32, data byte) {
	if address >= 0x100000 && address <= 0x10FFFF {
		baseAddress := uint32(0x100000)
		absoluteAddress := address - baseAddress
		s.ram[absoluteAddress] = data
	}
	if address >= 0x200000 && address <= 0x2fffff {
		baseAddress := uint32(0x200000)
		absoluteAddress := address - baseAddress
		s.flash[absoluteAddress] = data
	}
}
