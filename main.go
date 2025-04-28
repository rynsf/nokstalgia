package main

import (
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	dr "github.com/rynsf/nokstalgia/driver"
	ed "github.com/rynsf/nokstalgia/emulatedispatcher"
)

const (
	PW = 10
)

var keymap = map[int32]int{
	rl.KeyL: 7,
	rl.KeyU: 8,
	rl.KeyY: 9,
	rl.KeyN: 0xC,
	rl.KeyE: 0xA,
	rl.KeyI: 0xB,
	rl.KeyS: 0x1,
	rl.KeyT: 0x4,
}

func fillMem(data []byte, mem []byte, base int) int {
	for i := range data {
		mem[base+i] = data[i]
	}
	return base
}

func renderScreen() {
	W := int(dr.Locate("SCREEN_WIDTH"))
	H := int(dr.Locate("SCREEN_HEIGHT"))
	bgColor := rl.Color{R: 114, G: 164, B: 136, A: 255}
	fgColor := rl.Color{R: 13, G: 24, B: 20, A: 255}
	rl.BeginDrawing()
	rl.ClearBackground(bgColor)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			if screen[y][x] == 1 {
				rl.DrawRectangle(int32(x*PW), int32(y*PW), int32(PW), int32(PW), fgColor)
			}
		}
	}
	rl.EndDrawing()
}

var screen [][]int
var nok ed.CpuState

func memAlign(n uint32) uint32 {
	return n & ^uint32(0x3)
}

func main() {
	flashBin, err := os.ReadFile("./assets/flash.fls")
	if err != nil {
		panic(err)
	}
	var r [16]uint32
	ramlen := 0x40000
	ramBase := 0x100000
	ram := make([]byte, ramlen)
	flashBase := 0x200000
	flashLen := len(flashBin)
	dynRamLen := ramlen
	dynRamBase := ramBase + ramlen
	dynamicRam := make([]byte, dynRamLen)
	W := int(dr.Locate("SCREEN_WIDTH"))
	H := int(dr.Locate("SCREEN_HEIGHT"))
	screen = make([][]int, H)
	for i := 0; i < H; i++ {
		screen[i] = make([]int, W)
	}
	fillMem([]byte{0x3}, ram, 0x3FEA0)
	fillMem([]byte{0x1}, ram, 0x3FC4A)
	fillMem([]byte{0x1}, ram, 0x3FB65)
	// TODO: fill ram with ram init block
	fillMem([]byte{0x00, 0x32, 0x3C, 0xC4, 0x00, 0x32, 0x3B, 0x14,
		0x00, 0x32, 0x3B, 0x74, 0x00, 0x32, 0x3B, 0x44, 0x00, 0x32,
		0x3C, 0xDC, 0x00, 0x32, 0x3C, 0xDC, 0x00, 0x32, 0x3C, 0x1C,
		0x00, 0x32, 0x3B, 0xEC, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x32, 0x43, 0x64, 0x00, 0x32, 0x43, 0x7C,
		0x00, 0x32, 0x43, 0xAC, 0x00, 0x32, 0x43, 0xDC, 0x00, 0x32,
		0x48, 0x74, 0x00, 0x32, 0x48, 0x8C, 0x00, 0x32, 0x48, 0xA4,
		0x00, 0x32, 0x47, 0x24, 0x00, 0x32, 0x47, 0x54, 0x00, 0x32,
		0x47, 0x84, 0x00, 0x32, 0x44, 0xFC, 0x00, 0x32, 0x44, 0x3C,
		0x00, 0x32, 0x44, 0x9C, 0x00, 0x32, 0x3C, 0x04, 0x00, 0x32,
		0x45, 0x5C, 0x00, 0x32, 0x45, 0xBC, 0x00, 0x32, 0x46, 0x1C,
		0x00, 0x32, 0x46, 0xF4, 0x00, 0x32, 0x47, 0x0C, 0x00, 0x32,
		0x47, 0xB4, 0x00, 0x32, 0x47, 0xCC, 0x00, 0x32, 0x47, 0xFC,
		0x00, 0x32, 0x48, 0x2C, 0x00, 0x32, 0x48, 0x5C, 0x00, 0x32,
		0x46, 0x7C, 0x00, 0x32, 0x4D, 0x24, 0x00, 0x32, 0x4D, 0x84,
		0x00, 0x32, 0x4D, 0xE4, 0x00, 0x32, 0x4E, 0x44, 0x00, 0x32,
		0x4E, 0xA4, 0x00, 0x32, 0x4F, 0x04, 0x00, 0x32, 0x4C, 0xDC,
		0x00, 0x32, 0x4C, 0xF4, 0x00, 0x32, 0x4D, 0x0C}, ram, 0xDD80)

	nok = ed.InitCpu(r, false, false, false, false, ram, uint32(ramBase), uint32(ramlen), dynamicRam, uint32(dynRamBase), uint32(dynRamLen), flashBin, uint32(flashBase), uint32(flashLen))

	nok.ResetReg()
	nok.SetReg(memAlign(uint32(dynRamBase+dynRamLen-4)), 13) // SP
	nok.SetReg(0x0, 14)                                      // LR
	nok.SetReg(0x27C2A4, 15)                                 // PC
	nok.RunSubroutine()

	nok.ResetReg()
	nok.SetReg(0x5E2, 0)                                     // a1
	nok.SetReg(memAlign(uint32(dynRamBase+dynRamLen-4)), 13) // SP
	nok.SetReg(0x0, 14)                                      // LR
	nok.SetReg(0x28FD7C, 15)                                 // PC
	nok.RunSubroutine()

	nok.ResetReg()
	nok.SetReg(0x5DC, 0)                                     // a1
	nok.SetReg(memAlign(uint32(dynRamBase+dynRamLen-4)), 13) // SP
	nok.SetReg(0x0, 14)                                      // LR
	nok.SetReg(0x28FD7C, 15)                                 // PC
	nok.RunSubroutine()

	nok.ResetReg()
	nok.SetReg(0, 0)                                         // a1
	nok.SetReg(0, 1)                                         // a2
	nok.SetReg(memAlign(uint32(dynRamBase+dynRamLen-4)), 13) // SP
	nok.SetReg(0x0, 14)                                      // LR
	nok.SetReg(0x2E655C, 15)                                 // PC
	nok.RunSubroutine()

	nok.FillMem([]byte{0x2}, 0xD816)

	nok.ResetReg()
	nok.SetReg(1, 0)                                         // a1
	nok.SetReg(0, 1)                                         // a2
	nok.SetReg(memAlign(uint32(dynRamBase+dynRamLen-4)), 13) // SP
	nok.SetReg(0x0, 14)                                      // LR
	nok.SetReg(0x2E655C, 15)                                 // PC
	nok.RunSubroutine()

	nok.FillMem([]byte{0x0, 0xA}, 0xF6B8)
	nok.ResetReg()
	nok.SetReg(0x5AF, 0)                                     // a1
	nok.SetReg(memAlign(uint32(dynRamBase+dynRamLen-4)), 13) // SP
	nok.SetReg(0x0, 14)                                      // LR
	nok.SetReg(0x28E41E, 15)                                 // PC
	nok.RunSubroutine()

	nok.ResetReg()
	nok.SetReg(0x1144, 0)                                    // a1
	nok.SetReg(memAlign(uint32(dynRamBase+dynRamLen-4)), 13) // SP
	nok.SetReg(0x0, 14)                                      // LR
	nok.SetReg(0x28EEF8, 15)                                 // PC
	nok.RunSubroutine()

	nok.ResetReg()
	nok.SetReg(0x5DC, 0)                                     // a1
	nok.SetReg(memAlign(uint32(dynRamBase+dynRamLen-4)), 13) // SP
	nok.SetReg(0x0, 14)                                      // LR
	nok.SetReg(0x28E8E8, 15)                                 // PC
	nok.DumbState()
	nok.RunSubroutine()
	nok.DumbState()

	nok.SendToLcd(screen)

	rl.SetTraceLogLevel(rl.LogError)
	rl.InitWindow(int32(W*PW), int32(H*PW), "nokstalgia")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)
	for !rl.WindowShouldClose() {
		for k := range keymap {
			if rl.IsKeyPressed(k) {
				dr.Enq(0xC8, 1, [3]uint32{uint32(keymap[k])})
			}
			if rl.IsKeyReleased(k) {
				dr.Enq(0xC9, 1, [3]uint32{uint32(keymap[k])})
			}
		}
		dr.TimerTick()
		m, ok := dr.Deq()
		if ok {
			nok.ResetReg()
			id := m.GetId()
			argc := m.GetArgc()
			argv := m.GetArgv()
			nok.SetMessage(id, argc, argv)
			nok.SetReg(id, 0)
			nok.SetReg(memAlign(uint32(dynRamBase+dynRamLen-4)), 13) // SP
			nok.SetReg(0x0, 14)                                      // LR
			nok.SetReg(0x28E8E8, 15)                                 // PC
			nok.RunSubroutine()
			nok.SendToLcd(screen)
		}
		renderScreen()
	}
}
