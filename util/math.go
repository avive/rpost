package util

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/stretchr/testify/assert"
	"math"
	"math/big"
	"math/bits"
	"strings"
	"testing"
)

// POST math utils funcs

// Get big-endian bytes encoding of i
// Result is 1 to 8 bytes long. 0x0 bytes is returned for 0
func EncodeToBytes(i uint64) []byte {
	if i == 0 {
		return []byte{0x0}
	}

	iBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(iBuf, i)

	// bits needed to encode i
	l := bits.Len64(i)

	// bytes needed to encode i
	lb := l / 8
	if l%8 != 0 {
		lb += 1 // one more byte needed for the extra bits
	}

	return iBuf[8-lb:]
}

// Get the bool value of the nth bit of a value of a byte
// bit is defined from right to left so the LSB bit is at 0 and the MSB it is at 7.
// e.g. byte is bits at indexes [7|6|5|4|3|2|1|0] and 0x1 is 00000001, 0x2 is 00000010
func GetNthBit(b byte, bit uint64) bool {
	return b&byte(math.Pow(2, float64(bit))) != 0
}

// Decode a big-endian uint64 from its binary encoding of up to length bits
func Uint64Value(b bitarray.BitArray, length uint64) (uint64, error) {
	res := uint64(0)

	for i := uint64(0); i < length; i++ {
		bit, err := b.GetBit(i)
		if err != nil {
			return 0, err
		}
		res <<= 1
		if bit {
			res |= 1
		}
	}

	return res, nil
}

// Get string representation of a BitArray
func String(b bitarray.BitArray, size uint64) (string, error) {
	var sb strings.Builder
	for i := uint64(0); i < size; i++ {
		bit, err := b.GetBit(i)
		if err != nil {
			return "", err
		}

		if bit == false {
			sb.WriteString("0")
		} else {
			sb.WriteString("1")
		}
	}
	return sb.String(), nil
}

// Returns the max nonce that is l bits long.
// e.g. for l=256, nonce should be 2^256-1
func GetMaxNonce(l int) *big.Int {
	n := big.NewInt(0)
	for i := 0; i < l; i++ {
		n = n.SetBit(n, i, 1)
	}
	return n
}

// Returns an int representing l bits set to 1
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
// p := 1 / 2^l
func GetProbability(l uint) float64 {
	return 1.0 / math.Pow(2.0, float64(l))
}

// solve for l the equation: p := 1 / 2^ l
func GetDifficulty(p float64) uint {
	return uint(math.Ceil(math.Log2(1.0 / p)))
}

// clear the l msb bits of data considered as a big endian int
func clearMsbBits(l uint, data []byte) *big.Int {
	z := new(big.Int).SetBytes(data)
	firstBitIdx := len(data)*8 - int(l)
	lastBitIdx := len(data)*8 - 1
	for i := firstBitIdx; i <= lastBitIdx; i++ {
		z = z.SetBit(z, i, 0)
	}
	return z
}

// test helper - generate l random bytes
func Rnd(t *testing.T, l uint) []byte {
	res := make([]byte, l)
	_, err := rand.Read(res)
	assert.NoError(t, err)
	return res
}
