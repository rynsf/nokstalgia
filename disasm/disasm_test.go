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
