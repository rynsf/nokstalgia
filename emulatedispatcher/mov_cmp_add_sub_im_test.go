package emulatedispatcher

import (
	"math/rand"
	"testing"
	"time"

	"github.com/rynsf/nokstalgia/emulatedispatcher/gbaemu/gba"
)

func generateMovCmpAddSub(opcode, rd, offset8 uint16) uint16 {
	var inst uint16
	inst = inst | (0b1 << 13)
	inst = inst | (opcode << 11)
	inst = inst | (rd << 8)
	inst = inst | offset8
	return inst
}

func TestMovCmpAddSubMovInst(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rd := uint16(0); rd < uint16(8); rd++ {
		for offset8 := uint16(0); offset8 < uint16(256); offset8++ {
			inst = generateMovCmpAddSub(0, rd, offset8)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, MOV R%d, #%d, \ngot %v, \nwant %v", inst, rd, offset8, c, want)
			}
		}
	}
}

func TestMovCmpAddSubCmpInst(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rd := uint16(0); rd < uint16(8); rd++ {
		for offset8 := uint16(0); offset8 < uint16(256); offset8++ {
			inst = generateMovCmpAddSub(1, rd, offset8)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, CMP R%d, #%d, \ngot %v, \nwant %v", inst, rd, offset8, c, want)
			}
		}
	}
}

func TestMovCmpAddSubAddInst(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rd := uint16(0); rd < uint16(8); rd++ {
		for offset8 := uint16(0); offset8 < uint16(256); offset8++ {
			inst = generateMovCmpAddSub(2, rd, offset8)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, ADD R%d, #%d, \ngot %v, \nwant %v", inst, rd, offset8, c, want)
			}
		}
	}
}

func TestMovCmpAddSubSubInst(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rd := uint16(0); rd < uint16(8); rd++ {
		for offset8 := uint16(0); offset8 < uint16(256); offset8++ {
			inst = generateMovCmpAddSub(3, rd, offset8)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, SUB R%d, #%d, \ngot %v, \nwant %v", inst, rd, offset8, c, want)
			}
		}
	}
}
