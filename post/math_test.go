package post

import (
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
	testclearMSBBits(t)
}

func testclearMSBBits(t *testing.T) {
	mask := big.NewInt(0xffff)
	z := clearMsbBits(1, mask.Bytes())
	assert.Equal(t, uint64(32767), z.Uint64())

	z = clearMsbBits(2, mask.Bytes())
	assert.Equal(t, uint64(0x3fff), z.Uint64())

	mask = big.NewInt(0xffffffff)
	z = clearMsbBits(8, mask.Bytes())
	assert.Equal(t, uint64(0xffffff), z.Uint64())
}
