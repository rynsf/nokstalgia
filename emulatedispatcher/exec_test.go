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

func generateAddSub(iflag, op, rn, rs, rd uint16) uint16 {
	var inst uint16
	inst = inst | (0b11 << 11)
	inst = inst | (iflag << 10)
	inst = inst | (op << 9)
	inst = inst | (rn << 6)
	inst = inst | (rs << 3)
	inst = inst | rd
	return inst
}

func TestAddSubAddRegister(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rn := uint16(0); rn < uint16(8); rn++ {
		for rs := uint16(0); rs < uint16(8); rs++ {
			for rd := uint16(0); rd < uint16(8); rd++ {
				inst = generateAddSub(0, 0, rn, rs, rd)
				randomizeRegisters(&g, &c, rng)
				g.ThumbExec(inst)
				want = getCpuFromGba(&g)
				c.ExecInstruction(inst)

				if !pass(&c, &want) {
					t.Errorf("%X, ADD R%d, R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, rn, c, want)
				}
			}
		}
	}
}

func TestAddSubSubRegister(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rn := uint16(0); rn < uint16(8); rn++ {
		for rs := uint16(0); rs < uint16(8); rs++ {
			for rd := uint16(0); rd < uint16(8); rd++ {
				inst = generateAddSub(0, 1, rn, rs, rd)
				randomizeRegisters(&g, &c, rng)
				g.ThumbExec(inst)
				want = getCpuFromGba(&g)
				c.ExecInstruction(inst)

				if !pass(&c, &want) {
					t.Errorf("%X, SUB R%d, R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, rn, c, want)
				}
			}
		}
	}
}

func TestAddSubAddImmediate(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rn := uint16(0); rn < uint16(8); rn++ {
		for rs := uint16(0); rs < uint16(8); rs++ {
			for rd := uint16(0); rd < uint16(8); rd++ {
				inst = generateAddSub(1, 0, rn, rs, rd)
				randomizeRegisters(&g, &c, rng)
				g.ThumbExec(inst)
				want = getCpuFromGba(&g)
				c.ExecInstruction(inst)

				if !pass(&c, &want) {
					t.Errorf("%X, ADD R%d, R%d, #%d, \ngot %v, \nwant %v", inst, rd, rs, rn, c, want)
				}
			}
		}
	}
}

func TestAddSubSubImmediate(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rn := uint16(0); rn < uint16(8); rn++ {
		for rs := uint16(0); rs < uint16(8); rs++ {
			for rd := uint16(0); rd < uint16(8); rd++ {
				inst = generateAddSub(1, 1, rn, rs, rd)
				randomizeRegisters(&g, &c, rng)
				g.ThumbExec(inst)
				want = getCpuFromGba(&g)
				c.ExecInstruction(inst)

				if !pass(&c, &want) {
					t.Errorf("%X, SUB R%d, R%d, #%d, \ngot %v, \nwant %v", inst, rd, rs, rn, c, want)
				}
			}
		}
	}
}

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
