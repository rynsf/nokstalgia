package emulatedispatcher

func isAddCarry(r uint64) bool {
	return r > 0xffffffff // greater an 32 bit number
}

func isSubCarry(result uint64) bool {
	return result < 0x100000000
}

func isAddOverflow(lhs, rhs, result uint32) bool {
	v := ^(lhs ^ rhs) & (lhs ^ result) & 0x8000_0000
	return v != 0

}

func isSubOverflow(lhs, rhs, result uint32) bool {
	v := (lhs ^ rhs) & (lhs ^ result) & 0x8000_0000
	return v > 0
}
