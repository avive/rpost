package post

import (
	"github.com/avive/rpost/bstring"
	"github.com/avive/rpost/hashing"
	"math/big"
)

const (
	rootId = ""
)

type MerkleTreeWriter interface {
	Write() ([]byte, error)
}

type MerkleTreeReader interface {
	ReadProof(id Identifier) (MerkleProof, error) // Returns the path from a node identified by id to the root node
	ReadProofs(indices []*big.Int, n uint64) (MerkleProofs, error)
	Close() error
}

type Node struct {
	Id    Identifier
	Label Label
}

type MerkleProof []Node
type MerkleProofs []MerkleProof

type merkleTree struct {
	fileName string           // merkle tree data file full path and name
	l        uint             // number of bits in a post store entry
	n        uint             // table size T=2^n
	psr      StoreReader      // reader for post data
	h        hashing.HashFunc // Hx()
	f        bstring.BinaryStringFactory
	w        IStoreWriter // merkle tree store writer
	r        IStoreReader // merkle tree store reader
}

// n - merkle tree size = 2^n
func NewMerkleTreeReader(fileName string, l uint, n uint, h hashing.HashFunc) (MerkleTreeReader, error) {

	r, err := NewTreeStoreReader(fileName, n)
	if err != nil {
		return nil, err
	}

	res := &merkleTree{
		fileName, l, n, nil, h, bstring.NewSMBinaryStringFactory(), nil, r,
	}

	return res, nil
}

// n - store length. T = 2^n
func NewMerkleTreeWriter(psr StoreReader, fileName string, l uint, n uint,
	h hashing.HashFunc) (MerkleTreeWriter, error) {

	w, err := NewTreeStoreWriter(fileName, n-1)
	if err != nil {
		return nil, err
	}

	res := &merkleTree{
		fileName, l, n, psr, h, bstring.NewSMBinaryStringFactory(), w, nil,
	}

	return res, nil
}

// for each store table idx in indices, return the Merkle path from the node at that index to the merkle tree root
func (mt *merkleTree) ReadProofs(indices []*big.Int, n uint64) (MerkleProofs, error) {

	mps := make(MerkleProofs, len(indices))

	for idx, data := range indices {
		// n -> T = 2^n leaves
		// k - merkle tree height. Leaves: 2^(n-1) -> height := n -1

		bn, err := mt.f.NewBinaryStringFromInt(data.Uint64()>>1, uint(n-1))
		if err != nil {
			return nil, err
		}

		path, err := mt.ReadProof(Identifier(bn.GetStringValue()))
		if err != nil {
			return nil, err
		}

		mps[idx] = path
	}

	return mps, nil
}

// Returns the sibling nodes on the path from a node identified by id to the root
func (mt *merkleTree) ReadProof(id Identifier) (MerkleProof, error) {

	// fmt.Printf("Create merkle proof for node id: %s\n", id)
	res := make(MerkleProof, len(id))
	currNodeId, err := mt.f.NewBinaryString(string(id))
	if err != nil {
		return nil, err
	}

	// fmt.Printf("Getting siblings for %s\n", currNodeId.GetStringValue())

	nodeIds, err := currNodeId.GetBNSiblings(false)
	if err != nil {
		return nil, err
	}

	for i, nodeId := range nodeIds {

		idx := Identifier(nodeId.GetStringValue())

		// fmt.Printf("Reading merkle node for node id: %s\n", idx)

		l, err := mt.r.Read(idx)

		if err != nil {
			return nil, err
		}

		res[i] = Node{idx, l}
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

	comm, err := mt.write(rootId)
	if err != nil {
		return nil, err
	}

	err = mt.w.Close()
	if err != nil {
		return nil, err
	}

	return comm, nil
}

// visit a node identified by nodeId and returns its value
func (mt *merkleTree) write(nodeId string) ([]byte, error) {

	var leftNodeValue, rightNodeValue []byte

	if uint(len(nodeId)) == mt.n-1 {
		// Node is a merkle tree leaf
		// e.g. for n = 2 (post table size 4), node "0" and "1" of length 1 should be Merkle leafs
		// node is a Merkle leaf node - compute its value based on the data in the store
		// e.g. hash of left and right post table entries
		bs, err := mt.f.NewBinaryString(nodeId)
		if err != nil {
			return nil, err
		}

		// data index for left and right nodes
		idx := bs.GetValue() * 2
		leftNodeValue, err = mt.psr.ReadBytes(idx)
		if err != nil {
			return nil, err
		}

		rightNodeValue, err = mt.psr.ReadBytes(idx + 1)
		if err != nil {
			return nil, err
		}
	} else {
		// Node is an internal Merkle tree node
		// Recursively compute its value based on its children and store it
		var err error
		leftNodeValue, err = mt.write(nodeId + "0")
		if err != nil {
			return nil, err
		}

		rightNodeValue, err = mt.write(nodeId + "1")
		if err != nil {
			return nil, err
		}
	}

	digest := mt.h.Hash(leftNodeValue, rightNodeValue)
	mt.w.Write(Identifier(nodeId), digest)
	return digest, nil
}

// Close the reader if it is open
func (mt *merkleTree) Close() error {
	if mt.r != nil {
		err := mt.r.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
