package prover

import "github.com/avive/rpost/post"

type Proof struct {
	Nonces       []uint64
	MerkleProofs []post.MerkleProofs
}
