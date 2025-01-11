package emulatedispatcher

import (
	"math/rand"
	"testing"
	"time"

	"github.com/rynsf/nokstalgia/emulatedispatcher/gbaemu/gba"
)

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
