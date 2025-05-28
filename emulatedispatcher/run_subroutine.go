package emulatedispatcher

import (
	dr "github.com/rynsf/nokstalgia/driver"
)

func memAlign(n uint32) uint32 {
	return n & ^uint32(0x3)
}

func (s *CpuState) subroutineIsRunning() bool {
	return s.register[pc] != 0
}

func (s *CpuState) fillRam(src, dst, l uint32) {
	for i := uint32(0); i < l; i++ {
		s.write(dst+i, s.read(src+i))
	}
}

func (s *CpuState) fillRamInitBlock() {
	initBlockPtr := dr.Locate("RAM_INIT_BLOCK")
	l := s.read32(initBlockPtr)
	addr := s.read32(initBlockPtr + 4)
	for l != 0 {
		s.fillRam(initBlockPtr+8, addr, l)

		l = (((l + 3) >> 2) << 2)
		initBlockPtr += l + 8
		l = s.read32(initBlockPtr)
		addr = s.read32(initBlockPtr + 4)
	}
}

func InitCpu(r [16]uint32, v, c, z, n bool, ram []byte, ramBase, ramLen uint32, dynamicRam []byte, dynRamBase, dynRamLen uint32, flash []byte, flashBase, flashLen uint32) CpuState {
	dr.InitDynamicMem(0x4, 0xe000)
	specialFunc[dr.Locate("MALLOC")] = malloc
	specialFunc[dr.Locate("SEND_MESSAGE")] = send_message
	specialFunc[dr.Locate("BLOCK_ALLOC")] = block_alloc
	specialFunc[dr.Locate("BLOCK_DEALLOC")] = block_dealloc
	specialFunc[dr.Locate("BLOCK_ALLOC_NOWAIT")] = block_alloc_nowait
	specialFunc[dr.Locate("LOAD_GLOBAL_SETTINGS_VALUE")] = loadGlobalSettingsValue
	specialFunc[dr.Locate("READ_DIRECTORY_FILE")] = readDirectoryFile
	specialFunc[dr.Locate("DEV_FUNC_TRACE")] = doNothing
	specialFunc[dr.Locate("TI_ID_SEND")] = doNothing
	specialFunc[dr.Locate("OS_CONDITIONAL_INT_ENABLE")] = doNothing
	specialFunc[dr.Locate("OS_CONDITIONAL_INT_DISABLE")] = doNothing
	specialFunc[dr.Locate("DEV_DISP_REFRESH")] = doNothing
	specialFunc[dr.Locate("ENGINE_INITILIZE_FILE")] = doNothing
	specialFunc[dr.Locate("TRANSFER_DOWNLOADED_DATA_INFO")] = doNothing
	specialFunc[dr.Locate("WRITE_DIRECTORY_FILE")] = doNothing
	specialFunc[dr.Locate("FREE_DIRECTORY_FILE")] = doNothing
	specialFunc[dr.Locate("CREATE_MENU")] = doNothing
	specialFunc[dr.Locate("LOAD_SAVED_GAME")] = doNothing
	specialFunc[dr.Locate("TRANSLATE_UCS2")] = doNothing
	specialFunc[dr.Locate("OS_TIMER_STOP")] = os_timer_stop
	specialFunc[dr.Locate("OS_TIMER_START")] = os_timer_start
	specialFunc[dr.Locate("IND_CALL")] = doNothing
	specialFunc[dr.Locate("DEV_KEY_GAME_MODE_ENABLE")] = doNothing
	specialFunc[dr.Locate("TONE_CLASS_DISABLE")] = doNothing
	specialFunc[dr.Locate("GAME_LOAD_HIGHSCORE")] = gameLoadHighscore
	specialFunc[dr.Locate("SETTINGS_GET_VALUE")] = settingsGetValue
	specialFunc[dr.Locate("GAME_LINK_DISTANCE")] = gameLinkDistance

	s := CpuState{
		register:     r,
		sr:           cpsr{v, c, z, n},
		loc:          0,
		ram:          ram,
		ramBaseAdr:   ramBase,
		ramLen:       ramLen,
		dynamicRam:   dynamicRam,
		dynRamBase:   dynRamBase,
		dynRamLen:    dynRamLen,
		flash:        flash,
		flashBaseAdr: flashBase,
		flashLen:     flashLen,
	}
	s.fillRamInitBlock()
	return s
}

func (s *CpuState) RunSubroutine() {
	for s.subroutineIsRunning() {
		if s.specialFuncHandler() {
			continue
		}
		instruction := uint16(s.read16(s.register[pc]))
		if Debug {
			s.step()
		}
		s.loc = s.register[pc]
		s.register[pc] += 2
		s.execInstruction(instruction)
	}
}

func (s *CpuState) RunFunc(pc uint32, regs ...uint32) {
	if len(regs) > 13 {
		panic("runFunc: too many register values given")
	}
	s.ResetReg()
	for i, r := range regs {
		s.SetReg(r, i)
	}
	s.SetReg(memAlign(uint32(s.dynRamBase+s.dynRamLen)), 13)
	s.SetReg(0x0, 14)
	s.SetReg(pc, 15)
	s.RunSubroutine()
}
