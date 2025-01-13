package emulatedispatcher

import (
	"math/rand"

	"github.com/rynsf/nokstalgia/emulatedispatcher/gbaemu/gba"
)

func randomizeRegisters(g *gba.GBA, want *CpuState, rng *rand.Rand) {
	for i := 0; i < 16; i++ {
		r := rng.Uint32()
		g.Reg.R[i] = r
		want.register[i] = r
	}
	want.sr.carry = false
	want.sr.negative = false
	want.sr.zero = false
	want.sr.overflow = false
	g.SetCPSRFlag(28, false)
	g.SetCPSRFlag(29, false)
	g.SetCPSRFlag(30, false)
	g.SetCPSRFlag(31, false)
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
	for i := 0; i < 16; i++ {
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
