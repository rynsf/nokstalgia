package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDisassembler(t *testing.T) {
	var tests = []struct {
		codebuffer []byte
		pc         int
		want       string
	}{
		{[]byte{0x06, 0xEA}, 0, "LSL R2, R5, #27"},
	}
	for _, test := range tests {
		descr := fmt.Sprintf("disassembler(%d, %v)", test.pc, test.codebuffer)
		out = new(bytes.Buffer)
		disassembler(test.pc, test.codebuffer)
		got := out.(*bytes.Buffer).String()
		if got != test.want {
			t.Errorf("%s = %q, want %q", descr, out, test.want)
		}
	}
}

type Params struct {
	bytes  uint32
	bitLen int
}

func TestTwosCompliment(t *testing.T) {
	var tests = []struct {
		input Params
		want  int
	}{
		{Params{0xffff, 16}, -1},
		{Params{0x0000, 16}, 0},
		{Params{0x00f1, 8}, -15},
		{Params{0x0071, 8}, 113},
	}
	for _, test := range tests {
		if got := twosComplementToInt(test.input.bytes, test.input.bitLen); got != test.want {
			t.Errorf("twosComplementToInt(%04X, %d) = %d, want %d", test.input.bytes, test.input.bitLen, got, test.want)
		}
	}
}
