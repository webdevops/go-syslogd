package main

// Sets the bit at pos in the integer n
// Author: Kevin Burke
func setBit(n int, pos uint) int {
	n |= (1 << pos)
	return n
}

// Clears the bit at pos in n
// Author: Kevin Burke
func clearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

// check if number has bit
// Author: Kevin Burke
func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}
