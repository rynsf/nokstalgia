package main

import (
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	dr "github.com/rynsf/nokstalgia/driver"
	ed "github.com/rynsf/nokstalgia/emulatedispatcher"
)

const (
	pw = 20
	w  = 84
	h  = 48
)

var keymap = map[int32]int{
	rl.KeyL:      1,
	rl.KeyU:      2,
	rl.KeyY:      3,
	rl.KeyN:      4,
	rl.KeyE:      8,
	rl.KeyI:      6,
	rl.KeyM:      7,
	rl.KeyComma:  8,
	rl.KeyPeriod: 9,
}

func fillMem(data []byte, mem []byte, base int) int {
	for i := range data {
		mem[base+i] = data[i]
	}
	return base
}

func renderScreen() {
	bgColor := rl.Color{R: 114, G: 164, B: 136, A: 255}
	fgColor := rl.Color{R: 13, G: 24, B: 20, A: 255}
	rl.BeginDrawing()
	rl.ClearBackground(fgColor)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if screen[y][x] == 0 {
				rl.DrawRectangle(int32(x*pw), int32(y*pw), int32(pw), int32(pw), bgColor)
			}
		}
	}
	rl.EndDrawing()
}

var screen [][]int
var nok ed.CpuState

func main() {
	flashBin, err := os.ReadFile("./assets/snake.fls")
	if err != nil {
		panic(err)
	}
	var r [16]uint32
	ramlen := 0x10000
	ramBase := 0x100000
	ram := make([]byte, ramlen)
	flashBase := 0x200000
	flashLen := len(flashBin)
	dynRamLen := 0x10000
	dynRamBase := 0x110000
	dynamicRam := make([]byte, dynRamLen)
	fillMem([]byte{0x01}, ram, 0xFCF0) // current language
	fillMem([]byte{0x03}, ram, 0xFFA8) // device power status
	fillMem([]byte{0x00, 0x00, 0x00, 0x01}, ram, 0x73D0)
	fillMem([]byte{0x00, 0x00, 0x00, 0x01}, ram, 0x73D4)
	fillMem([]byte{0x00, 0x00, 0x00, 0x56}, ram, 0xB6A0) // set the seed for rand
	screen = make([][]int, 48)
	for i := 0; i < 48; i++ {
		screen[i] = make([]int, 84)
	}

	rl.SetTraceLogLevel(rl.LogError)
	rl.InitWindow(int32(w*pw), int32(h*pw), "nokstalgia")
	defer rl.CloseWindow()

	nok = ed.InitCpu(r, false, false, false, false, ram, uint32(ramBase), uint32(ramlen), dynamicRam, uint32(dynRamBase), uint32(dynRamLen), flashBin, uint32(flashBase), uint32(flashLen))

	// Init game by calling the menu dispatcher with the start_game message
	nok.SetReg(0x1C00, 0)
	nok.SetReg(0x11fff8, 13)
	nok.SetReg(0x0, 14)
	nok.SetReg(0x2AAD68, 15)
	nok.RunSubroutine()

	setSubroutine := func() {
		nok.SetReg(0x2AB118, 15) // game dispatcher
		nok.SetReg(0x0, 14)
		nok.SetReg(0x11fff8, 13)
	}

	dr.Enq(0x5DC, 0, [3]uint32{}) // Enqueue message_d_init, for the game dispatcher

	rl.SetTargetFPS(120)
	for !rl.WindowShouldClose() {
		// handle key events
		for k := range keymap {
			if rl.IsKeyPressed(k) {
				dr.Enq(0xC8, 1, [3]uint32{uint32(keymap[k])})
			}
		}
		dr.TimerTick()
		m, ok := dr.Deq()
		if ok {
			setSubroutine()
			id := m.GetId()
			argc := m.GetArgc()
			argv := m.GetArgv()
			nok.SetMessage(id, argc, argv)
			nok.SetReg(id, 0)
			nok.RunSubroutine()
			nok.UpdateScreen()
			nok.SendToLcd(screen)
		}
		renderScreen()
		if nok.GetReg(0) == 0x5DF {
			break
		}
	}
}
