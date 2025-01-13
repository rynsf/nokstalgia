package emulatedispatcher

import (
	"math/rand"
	"testing"
	"time"

	"github.com/rynsf/nokstalgia/emulatedispatcher/gbaemu/gba"
)

func generateHiRegOpBranchEx(op, h1, h2, rs, rd uint16) uint16 {
	var inst uint16
	inst = inst | (1 << 14)
	inst = inst | (1 << 10)
	inst = inst | (op << 8)
	inst = inst | (h1 << 7)
	inst = inst | (h2 << 6)
	inst = inst | (rs << 3)
	inst = inst | rd
	return inst
}

func TestHiRegOpBranchExAdd(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for h1 := uint16(0); h1 < 2; h1++ {
		for h2 := uint16(0); h2 < 2; h2++ {
			for rs := uint16(0); rs < uint16(8); rs++ {
				for rd := uint16(0); rd < uint16(8); rd++ {
					inst = generateHiRegOpBranchEx(0, h1, h2, rs, rd)
					randomizeRegisters(&g, &c, rng)
					g.ThumbExec(inst)
					want = getCpuFromGba(&g)
					c.ExecInstruction(inst)

					if !pass(&c, &want) {
						t.Errorf("%X, ADD R%d, R%d, \ngot %v, \nwant %v", inst, h1*8+rd, h2*8+rs, c, want)
					}
				}
			}
		}
	}
}

func TestHiRegOpBranchExCmp(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for h1 := uint16(0); h1 < 2; h1++ {
		for h2 := uint16(0); h2 < 2; h2++ {
			for rs := uint16(0); rs < uint16(8); rs++ {
				for rd := uint16(0); rd < uint16(8); rd++ {
					inst = generateHiRegOpBranchEx(1, h1, h2, rs, rd)
					randomizeRegisters(&g, &c, rng)
					g.ThumbExec(inst)
					want = getCpuFromGba(&g)
					c.ExecInstruction(inst)

					if !pass(&c, &want) {
						t.Errorf("%X, CMP R%d, R%d, \ngot %v, \nwant %v", inst, h1*8+rd, h2*8+rs, c, want)
					}
				}
			}
		}
	}
}

func TestHiRegOpBranchExMov(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for h1 := uint16(0); h1 < 2; h1++ {
		for h2 := uint16(0); h2 < 2; h2++ {
			for rs := uint16(0); rs < uint16(8); rs++ {
				for rd := uint16(0); rd < uint16(8); rd++ {
					inst = generateHiRegOpBranchEx(2, h1, h2, rs, rd)
					randomizeRegisters(&g, &c, rng)
					g.ThumbExec(inst)
					want = getCpuFromGba(&g)
					c.ExecInstruction(inst)

					if !pass(&c, &want) {
						t.Errorf("%X, MOV R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
					}
				}
			}
		}
	}
}

func TestHiRegOpBranchExBx(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for h1 := uint16(0); h1 < 2; h1++ {
		for h2 := uint16(0); h2 < 2; h2++ {
			for rs := uint16(0); rs < uint16(8); rs++ {
				for rd := uint16(0); rd < uint16(8); rd++ {
					inst = generateHiRegOpBranchEx(3, h1, h2, rs, rd)
					randomizeRegisters(&g, &c, rng)
					g.ThumbExec(inst)
					want = getCpuFromGba(&g)
					c.ExecInstruction(inst)

					if !pass(&c, &want) {
						t.Errorf("%X, BX R%d, \ngot %v, \nwant %v", inst, rs, c, want)
					}
				}
			}
		}
	}
}
