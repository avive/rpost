package post

import (
	"fmt"
	"math/big"
)

// POST math utils funcs

// Get the prob (0...1) of a difficulty param l
func difficultyToProb(l uint) float64 {
	r := 100.0
	for i := uint(0); i < l; i++ {
		r = r / 2
	}
	return r / 100
}

// clear the l msb bits of data considered as a big endian int
func clearMSBBits(l uint, data []byte) *big.Int {
	z := new(big.Int).SetBytes(data)
	firstBitIdx := len(data)*8 - int(l)
	lastBitIdx := len(data)*8 - 1

	for i := firstBitIdx; i <= lastBitIdx; i++ {
		z = z.SetBit(z, i, 0)
	}

	fmt.Printf("input mask:\n %s \n 0x%x \n", new(big.Int).SetBytes(data).String(), data)
	fmt.Printf("new mask:\n %s \n 0x%x \n", z.String(), z.Bytes())
	return z
}
