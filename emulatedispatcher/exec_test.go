package emulatedispatcher

import (
	"math/rand"
	"testing"
	"time"

	"github.com/rynsf/nokstalgia/emulatedispatcher/gbaemu/gba"
)

func randomizeRegisters(g *gba.GBA, want *CpuState, rng *rand.Rand) {
	for i := 0; i < 12; i++ {
		r := rng.Uint32()
		g.Reg.R[i] = r
		want.register[i] = r
	}
}

func generateMovShiftRegLSL(offset, rs, rd uint16) uint16 {
	var inst uint16
	inst = inst | (offset << 6) // set offset
	inst = inst | (rs << 3)     // set rs
	inst = inst | rd            // set rd
	return inst
}

func getCpuFromGba(g *gba.GBA) CpuState {
	c := CpuState{}
	for i := 0; i < 16; i++ {
		c.register[i] = g.R[i]
	}
	c.sr.carry = g.GetCPSRFlag(gba.FlagC)
	c.sr.negative = g.GetCPSRFlag(gba.FlagN)
	c.sr.overflow = g.GetCPSRFlag(gba.FlagV)
	c.sr.zero = g.GetCPSRFlag(gba.FlagZ)
	return c
}

func pass(c, want *CpuState) bool {
	for i := 0; i < 12; i++ {
		if c.register[i] != want.register[i] {
			return false
		}
	}
	if c.sr.carry != want.sr.carry {
		return false
	}
	if c.sr.negative != want.sr.negative {
		return false
	}
	if c.sr.zero != want.sr.zero {
		return false
	}
	if c.sr.overflow != want.sr.overflow {
		return false
	}
	return true
}

func TestMovShiftRegLSL(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for offset := uint16(0); offset < uint16(32); offset++ {
		for rs := uint16(0); rs < uint16(8); rs++ {
			for rd := uint16(0); rd < uint16(8); rd++ {
				inst = generateMovShiftRegLSL(offset, rs, rd)
				randomizeRegisters(&g, &c, rng)
				g.ThumbExec(inst)
				want = getCpuFromGba(&g)
				c.ExecInstruction(inst)

				if !pass(&c, &want) {
					t.Errorf("LSL R%d, R%d, #%d, \ngot %v, \nwant %v", rd, rs, offset, c, want)
				}
			}
		}
	}
}

func generateMovShiftRegLSR(offset, rs, rd uint16) uint16 {
	var inst uint16
	inst = inst | (0b01 << 11)  // set opcode
	inst = inst | (offset << 6) // set offset
	inst = inst | (rs << 3)     // set rs
	inst = inst | rd            // set rd
	return inst
}

func TestMovShiftRegLSR(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for offset := uint16(0); offset < uint16(32); offset++ {
		for rs := uint16(0); rs < uint16(8); rs++ {
			for rd := uint16(0); rd < uint16(8); rd++ {
				inst = generateMovShiftRegLSR(offset, rs, rd)
				randomizeRegisters(&g, &c, rng)
				g.ThumbExec(inst)
				want = getCpuFromGba(&g)
				c.ExecInstruction(inst)

				if !pass(&c, &want) {
					t.Errorf("LSR R%d, R%d, #%d, \ngot %v, \nwant %v", rd, rs, offset, c, want)
				}
			}
		}
	}
}

func generateMovShiftRegASR(offset, rs, rd uint16) uint16 {
	var inst uint16
	inst = inst | (0b10 << 11)  // set opcode
	inst = inst | (offset << 6) // set offset
	inst = inst | (rs << 3)     // set rs
	inst = inst | rd            // set rd
	return inst
}

func TestMovShiftRegASR(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for offset := uint16(0); offset < uint16(32); offset++ {
		for rs := uint16(0); rs < uint16(8); rs++ {
			for rd := uint16(0); rd < uint16(8); rd++ {
				inst = generateMovShiftRegASR(offset, rs, rd)
				randomizeRegisters(&g, &c, rng)
				g.ThumbExec(inst)
				want = getCpuFromGba(&g)
				c.ExecInstruction(inst)

				if !pass(&c, &want) {
					t.Errorf("%X, ASR R%d, R%d, #%d, \ngot %v, \nwant %v", inst, rd, rs, offset, c, want)
				}
			}
		}
	}
}
