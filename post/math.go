package post

import (
	"math/big"
)

// POST math utils funcs

// Returns a an int representing l 1 bits.
func GetSimpleMask(l uint) *big.Int {
	mask := big.NewInt(0)
	for x := 0; x < int(l); x++ {
		mask = mask.SetBit(mask, x, 1)
	}
	return mask
}

// Create an l bytes long bit mask with the c MSB bits set to 0 and the other bits set to 1
// c*8 must be <= l
func GetMask(l uint, c uint) *big.Int {

	// Generate a mask here of 32 bytes with l 0 bits at msb
	mask := make([]byte, l)
	for i := uint(0); i < l; i++ {
		mask[i] = 0xff
	}

	return clearMsbBits(c, mask)
}

// Get the prob (0...1) of a difficulty param l
func GetProbability(l uint) float64 {
	r := 100.0
	for i := uint(0); i < l; i++ {
		r = r / 2
	}
	return r / 100
}

// clear the l msb bits of data considered as a big endian int
func clearMsbBits(l uint, data []byte) *big.Int {
	z := new(big.Int).SetBytes(data)
	firstBitIdx := len(data)*8 - int(l)
	lastBitIdx := len(data)*8 - 1

	for i := firstBitIdx; i <= lastBitIdx; i++ {
		z = z.SetBit(z, i, 0)
	}

	// fmt.Printf("input mask:\n %s \n 0x%x \n", new(big.Int).SetBytes(data).String(), data)
	// fmt.Printf("new mask:\n %s \n 0x%x \n", z.String(), z.Bytes())
	return z
}
