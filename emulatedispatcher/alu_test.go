package emulatedispatcher

import (
	"math/rand"
	"testing"
	"time"

	"github.com/rynsf/nokstalgia/emulatedispatcher/gbaemu/gba"
)

const (
	and = iota
	eor
	lsl
	lsr
	asr
	adc
	sbc
	ror
	tst
	neg
	cmp
	cmn
	orr
	mul
	bic
	mvn
)

func generateAlu(op, rs, rd uint16) uint16 {
	var inst uint16
	inst = inst | (0b1 << 14)
	inst = inst | (op << 6)
	inst = inst | (rs << 3)
	inst = inst | rd
	return inst
}

func TestAluAnd(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(and, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, AND R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluEor(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(eor, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, EOR R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluLsl(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(lsl, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, LSL R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluLsr(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(lsr, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, LSR R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluAsr(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(asr, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, ASR R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluAdc(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(adc, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, ADC R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluSbc(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(sbc, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, SBC R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluRor(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(ror, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, ROR R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluTst(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(tst, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, TST R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluNeg(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(neg, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, NEG R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluCmp(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(cmp, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, CMP R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluCmn(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(cmn, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, CMN R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluOrr(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(orr, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, ORR R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluMul(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(mul, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, MUL R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluBic(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(bic, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, BIC R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}

func TestAluMvn(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random Seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	r := gba.Reg{}
	g := gba.GBA{Reg: r}
	var inst uint16
	var c CpuState
	var want CpuState

	for rs := uint16(0); rs < uint16(8); rs++ {
		for rd := uint16(0); rd < uint16(8); rd++ {
			inst = generateAlu(mvn, rs, rd)
			randomizeRegisters(&g, &c, rng)
			g.ThumbExec(inst)
			want = getCpuFromGba(&g)
			c.ExecInstruction(inst)

			if !pass(&c, &want) {
				t.Errorf("%X, MVN R%d, R%d, \ngot %v, \nwant %v", inst, rd, rs, c, want)
			}
		}
	}
}
