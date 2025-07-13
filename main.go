package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	dr "github.com/rynsf/nokstalgia/driver"
	ed "github.com/rynsf/nokstalgia/emulatedispatcher"
)

const (
	PW = 10
)

// TODO: add custom per game key bindings
/*
bindings for colemak wide mode
var keymap = map[int32]int{
	rl.KeyEight:  0x1,
	rl.KeyNine:   0x2,
	rl.KeyZero:   0x3,
	rl.KeyL:      0x4,
	rl.KeyU:      0x5,
	rl.KeyY:      0x6,
	rl.KeyN:      0x7,
	rl.KeyE:      0x8,
	rl.KeyI:      0x9,
	rl.KeyM:      0xC,
	rl.KeyComma:  0xA,
	rl.KeyPeriod: 0xB,
	rl.KeyW:      0x19,
	rl.KeyF:      0x17,
	rl.KeyP:      0x1A,
	rl.KeyR:      0xE,
	rl.KeyS:      0x18,
	rl.KeyT:      0xF,
}
*/

// bindings for qwerty
var keymap = map[int32]int{
	rl.KeySeven:  0x1,
	rl.KeyEight:  0x2,
	rl.KeyNine:   0x3,
	rl.KeyU:      0x4,
	rl.KeyI:      0x5,
	rl.KeyO:      0x6,
	rl.KeyJ:      0x7,
	rl.KeyK:      0x8,
	rl.KeyL:      0x9,
	rl.KeyM:      0xC,
	rl.KeyComma:  0xA,
	rl.KeyPeriod: 0xB,
	rl.KeyW:      0x19,
	rl.KeyE:      0x17,
	rl.KeyR:      0x1A,
	rl.KeyS:      0xE,
	rl.KeyD:      0x18,
	rl.KeyF:      0xF,
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

var gameList = []string{ //TODO: write locate to find all games in rom
	"Snake II",
	"Space Impact",
	"Link5",
}

func printGameList() {
	fmt.Println("Games available in ROM are:")
	for i, s := range gameList {
		fmt.Printf("%d. %s\n", i+1, s)
	}
}

func usage() {
	fmt.Println("Welcome to Nokstalgia!") //TODO: Write a friendly manual
}

func whichGame(game string) string {
	match := []string{}
	for _, s := range gameList {
		if strings.Contains(strings.ToLower(s), strings.ToLower(game)) {
			match = append(match, s)
		}
	}
	if len(match) != 1 {
		return ""
	}
	return match[0]
}

func parseArgs() string {
	if len(os.Args) < 3 {
		usage()
		return ""
	}
	switch os.Args[1] {
	case "run":
		runCmdSet := flag.NewFlagSet("run", flag.ExitOnError)
		if len(os.Args) < 4 {
			usage()
			return ""
		}
		selectedGame := whichGame(os.Args[2])
		if selectedGame == "" {
			return ""
		}
		switch selectedGame {
		case "Snake II":
			level := runCmdSet.Int("level", 0, "Speed of snake")
			maze := runCmdSet.Int("maze", 0, "maze")
			runCmdSet.Parse(os.Args[3:])
			dr.SetConfig("level", uint32(*level))
			dr.SetConfig("maze", uint32(*maze))
			dr.SetSelectedGame("Snake II")
			return runCmdSet.Args()[0]
		case "Space Impact":
			dr.SetSelectedGame("Space Impact")
			return os.Args[len(os.Args)-1]
		case "Link5":
			level := runCmdSet.Int("level", 0, "level")
			rules := runCmdSet.Int("rules", 0, "rules")
			challenges := runCmdSet.Int("challenges", 0, "challenges")
			runCmdSet.Parse(os.Args[3:])
			dr.SetConfig("level", uint32(*level))
			dr.SetConfig("rules", uint32(*rules))
			dr.SetConfig("challenges", uint32(*challenges))
			dr.SetSelectedGame("Link5")
			return os.Args[len(os.Args)-1]
		}
	case "list":
		printGameList()
		return ""
	default:
		usage()
		return ""
	}
	return ""
}

func main() {
	romName := parseArgs()
	if romName == "" {
		return
	}
	flashBin, err := os.ReadFile(romName)
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

	switch dr.GetSelectedGame() {
	case "Snake II":
		initGameSnake(nok)
	case "Space Impact":
		initGameSpace(nok)
	case "Link5":
		initGameLink(nok)
	}

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
