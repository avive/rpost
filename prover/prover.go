package prover

import (
	"github.com/avive/rpost/hashing"
	"github.com/avive/rpost/post"
)

const (
	k = 256 // this can't be modified and is set as it needs to be the same as the bit length of the output of sha256()
)

type Table struct {
	id []byte           // initial commitment
	n  uint64           // n param 1 <= n <= 63 - table size is 2^n
	l  uint             // l param (num of leading 0s for p) := f(p). 1: 50%, 2: 25%, 3:12.5%...
	h  hashing.HashFunc // Hx()
	sr post.StoreReader
}
