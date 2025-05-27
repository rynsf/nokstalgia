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

func initGameSnake(nok ed.CpuState) {
	nok.FillMem([]byte{0, 0, 0, 0x19}, 0xFC68)
	nok.RunFunc(0x28FD7C, 0x5DC)
	nok.RunFunc(0x2E655C, 0, 0)

	nok.FillMem([]byte{0x1}, 0xD816)
	nok.RunFunc(0x2E655C, 0, 0)

	for i := 0; i < 9; i++ {
		nok.RunFunc(0x28E41E, 0x5AF)
	}

	nok.RunFunc(0x28EEF8, 0x1144)
	nok.RunFunc(0x28E8E8, 0x5DC)
}

func initGameLink(nok ed.CpuState) {
	nok.RunFunc(0x28FD7C, 0x5DC)
	nok.RunFunc(0x2E655C, 0, 0)

	nok.FillMem([]byte{0x3}, 0xD816)
	nok.RunFunc(0x2E655C, 2, 0)

	nok.FillMem([]byte{0x7}, 0xFB5C)
	nok.RunFunc(0x28E41E, 0x5AF)

	nok.RunFunc(0x28EEF8, 0x1144)
	nok.RunFunc(0x28E8E8, 0x5DC)
}

func initGameSpace(nok ed.CpuState) {
	nok.RunFunc(0x28FD7C, 0x5DC)
	nok.RunFunc(0x2E655C, 0, 0)

	nok.FillMem([]byte{0x2}, 0xD816)
	nok.RunFunc(0x2E655C, 1, 0)

	nok.FillMem([]byte{0x0, 0xA}, 0xF6B8)
	nok.RunFunc(0x28E41E, 0x5AF)

	nok.RunFunc(0x28EEF8, 0x1144)
	nok.RunFunc(0x28E8E8, 0x5DC)
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

	nok = ed.InitCpu(r, false, false, false, false, ram, uint32(ramBase), uint32(ramlen), dynamicRam, uint32(dynRamBase), uint32(dynRamLen), flashBin, uint32(flashBase), uint32(flashLen))

	nok.RunFunc(0x27C2A4)
	nok.RunFunc(0x28FD7C, 0x5E2)

	initGameSpace(nok)

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
			id := m.GetId()
			nok.SetMessage(id, m.GetArgc(), m.GetArgv())
			nok.RunFunc(0x28E8E8, id)
			nok.SendToLcd(screen)
		}
		renderScreen()
	}
}
