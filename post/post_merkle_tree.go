package post

import (
	"github.com/avive/rpost/hashing"
	"github.com/avive/rpost/merkle"
)

const (
	rootId = ""
)

type IMerkleTreeWriter interface {
	Write() ([]byte, error)
}

type merkleTree struct {
	fileName string
	l        uint             // number of bits in a post store entry
	n        uint             // table size T=2^n
	psr      StoreReader      // reader for post data
	h        hashing.HashFunc // Hx()
	f        merkle.BinaryStringFactory
	w        merkle.IStoreWriter // merkle tree store writer
}

func NewMerkleTreeWriter(postStore string, fileName string, l uint, n uint,
	h hashing.HashFunc) (IMerkleTreeWriter, error) {

	psr, err := NewStoreReader(postStore, l)
	if err != nil {
		return nil, err
	}

	w, err := merkle.NewTreeStoreWriter(fileName, n-1)
	if err != nil {
		return nil, err
	}

	res := &merkleTree{
		fileName, l, n, psr, h, merkle.NewSMBinaryStringFactory(), w,
	}

	return res, nil
}

// Write the Merkle tree of the provided store to the store
// Returns the Merkle root commitment for the data
func (mt *merkleTree) Write() ([]byte, error) {

	// Merkle tree height equals to log of data size minus 1
	// h := mt.n - 1
	// Number of table entries
	// t := uint64(math.Pow(2, float64(mt.n)))

	return mt.visit(rootId)
}

// visit a node identified by nodeId and returns its value
func (mt *merkleTree) visit(nodeId string) ([]byte, error) {

	if uint(len(nodeId)) == mt.n-1 {
		// Node is a merkle tree leaf
		// e.g. for n = 2 (post table size 4), node "0" and "1" of length 1 should be Merkle leafs
		// node is a Merkle leaf node - compute its value based on the data in the store
		// e.g. hash of left and right post table entries

		bs, err := mt.f.NewBinaryString(nodeId)
		if err != nil {
			return nil, err
		}

		idx := bs.GetValue() * 2
		leftValue, err := mt.psr.ReadBytes(idx - 1)
		if err != nil {
			return nil, err
		}

		rightValue, err := mt.psr.ReadBytes(idx)
		if err != nil {
			return nil, err
		}

		digest := mt.h.Hash(leftValue, rightValue)
		mt.w.Write(merkle.Identifier(nodeId), digest)

		return digest, nil
	}

	// Node is an internal Merkle tree node
	// Recursively compute its value based on its children and store it

	leftNodeId := nodeId + "0"
	rightNodeId := nodeId + "1"

	leftNodeValue, err := mt.visit(leftNodeId)
	if err != nil {
		return nil, err
	}

	rightNodeValue, err := mt.visit(rightNodeId)
	if err != nil {
		return nil, err
	}

	digest := mt.h.Hash(leftNodeValue, rightNodeValue)
	mt.w.Write(merkle.Identifier(nodeId), digest)

	return digest, nil

}
