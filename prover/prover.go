package prover

import (
	"github.com/avive/rpost/hashing"
	"github.com/avive/rpost/post"
)

type Prover interface {
	Prove(challenge []byte) error
}

// todo: define proof

type prover struct {
	id []byte           // initial commitment
	n  uint64           // n param 1 <= n <= 63 - table size is 2^n
	l  uint             // l param (num of leading 0s for p) := f(p). 1: 50%, 2: 25%, 3:12.5%...
	h  hashing.HashFunc // Hx()
	sr post.StoreReader
	mr post.MerkleTreeReader
}

func NewProver(id []byte, n uint64, l uint, h hashing.HashFunc, storeFile string, merkleFile string) (Prover, error) {

	sr, err := post.NewStoreReader(storeFile, uint(n))
	if err != nil {
		return nil, err
	}

	mr, err := post.NewMerkleTreeReader(merkleFile, l, uint(n-1), h)
	if err != nil {
		return nil, err
	}

	prover := &prover{
		id, n, l, h, sr, mr,
	}

	return prover, nil
}

func (p *prover) Prove(challenge []byte) error {
	return nil
}
