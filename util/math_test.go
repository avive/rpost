package util

import (
	"fmt"
	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestIntMask(t *testing.T) {
	i := GetSimpleMask(20)
	e := big.NewInt(1048575)
	assert.Equal(t, e, i)

	i = GetSimpleMask(32)
	e = big.NewInt(4294967295)
	assert.Equal(t, e, i)
}

func TestZeroLsbsMask(t *testing.T) {
	m := GetMask(32, 20)
	s := "110427941548649020598956093796432407239217743554726184882600387580788735"
	assert.Equal(t, s, m.String())

	m = GetMask(1, 8)
	assert.Equal(t, uint64(0x0), m.Uint64())

	m = GetMask(1, 4)
	assert.Equal(t, uint64(0xf), m.Uint64())

	m = GetMask(2, 0)
	assert.Equal(t, uint64(0xffff), m.Uint64())

	m = GetMask(2, 1)
	assert.Equal(t, uint64(0x7fff), m.Uint64())

	m = GetMask(2, 2)
	assert.Equal(t, uint64(0x3fff), m.Uint64())

}

func TestGetProbability(t *testing.T) {

	p := GetProbability(0)
	assert.Equal(t, 1.0, p)

	p = GetProbability(1)
	assert.Equal(t, 0.5, p)

	p = GetProbability(2)
	assert.Equal(t, 0.25, p)

	p = GetProbability(3)
	assert.Equal(t, 0.125, p)

	p = GetProbability(4)
	assert.Equal(t, 0.0625, p)
}

func TestAll(t *testing.T) {
	testClearMsbBits(t)
}

func testClearMsbBits(t *testing.T) {
	mask := big.NewInt(0xffff)
	z := clearMsbBits(1, mask.Bytes())
	assert.Equal(t, uint64(32767), z.Uint64())

	z = clearMsbBits(2, mask.Bytes())
	assert.Equal(t, uint64(0x3fff), z.Uint64())

	mask = big.NewInt(0xffffffff)
	z = clearMsbBits(8, mask.Bytes())
	assert.Equal(t, uint64(0xffffff), z.Uint64())
}

func TestGetMaxNonce(t *testing.T) {
	n := GetMaxNonce(256)
	s := fmt.Sprintf("%x", n)
	assert.Equal(t, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", s)
	// fmt.Printf("%x", n)
}

func TestBitArrayString(t *testing.T) {
	b := bitarray.NewBitArray(20, false)
	err := b.SetBit(0)
	assert.NoError(t, err)

	err = b.SetBit(18)
	assert.NoError(t, err)

	s, err := String(b, 20)
	assert.NoError(t, err)
	fmt.Printf(s)
}

func TestGetNthBit(t *testing.T) {
	b := byte(1)
	assert.True(t, GetNthBit(b, 0), "expected bit 0 to be set for 0x1")
	assert.False(t, GetNthBit(b, 1), "expected bit 1 to be 0 for 0x1")

	b = byte(2)
	assert.True(t, GetNthBit(b, 1))
	assert.False(t, GetNthBit(b, 0))

}
