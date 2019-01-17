package prover

import (
	"errors"
	"fmt"
	"github.com/avive/rpost/hashing"
	"github.com/avive/rpost/post"
	"github.com/avive/rpost/util"
	"math"
	"math/big"
)

const K = post.K

type Prover interface {
	Prove(challenge []byte) (*Proof, error)
}

type Proof struct {
	nonces []uint64
	mpaths []post.MerklePaths
}

// todo: define proof tpye

type prover struct {
	id []byte                // initial commitment
	n  uint64                // n param 9 <= n <= 63 - table size is 2^n
	l  uint                  // l param (num of leading 0s for p) := f(p). 1: 50%, 2: 25%, 3:12.5%...
	h  hashing.HashFunc      // Hx()
	sr post.StoreReader      // Store reader can read data from the store at any index
	mr post.MerkleTreeReader // Merkle tree reader can read nodes on the path from an identified nodes the root
}

func NewProver(id []byte, n uint64, l uint, h hashing.HashFunc, storeFile string, merkleFile string) (Prover, error) {

	if n < 9 {
		return nil, errors.New("n must be >= 9")
	}

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

func (p *prover) Prove(challenge []byte) (*Proof, error) {

	// implements the prover proof phase described in page 9 of the paper

	// table size as big int
	T := big.NewInt(int64(math.Pow(2, float64(p.n))))

	// holds i(j,t) indexes as defined in page 9
	indices := make([]*big.Int, K)

	// holds nonce(j)
	nonces := make([]uint64, K)

	temp := new(big.Int)

	// hold K merkle paths. e.g. Phi(decommit(i))
	mpaths := make([]post.MerklePaths, K)

	var mpj post.MerklePaths

	// compute big int mask for pathProbe < phi calculations
	phi := float64(K) / float64(T.Uint64())
	diff := util.GetDifficulty(phi)
	mask := util.GetMask(32, diff)
	fmt.Printf("Mask : %s\n", mask.String())

	for j := 0; j < K; j++ {
		nonce := uint64(0)
		for {
			nonce += 1
			for t := 0; t < K; t++ {
				nb := util.EncodeToBytes(nonce)
				d := p.h.Hash(nb, p.id, []byte{byte(j)}, []byte{byte(t)})
				temp = temp.SetBytes(d)
				indices[t] = temp.Mod(temp, T)
			}

			// read merkle paths
			mpj = p.mr.ReadMerklePaths(indices)
			pathProbe, err := p.computePathProbe(indices, mpj)
			if err != nil {
				return nil, err
			}

			if pathProbe.Cmp(mask) == -1 {
				break
			}
		}

		mpaths[j] = mpj
		nonces[j] = nonce
	}

	return &Proof{nonces, mpaths}, nil
}

func (p *prover) computePathProbe(indices []*big.Int, mpj post.MerklePaths) (*big.Int, error) {

	data := make([][]byte, len(indices)*2+len(mpj))

	for i := 0; i < K; i++ {
		data[i*2] = indices[i].Bytes()

		// read the data from the store
		buff, err := p.sr.ReadBytes(indices[i].Uint64())
		if err != nil {
			return nil, err
		}
		data[i*2+1] = buff
	}

	idx := K*2 + 2

	var buf []byte
	for _, path := range mpj {
		for _, node := range path {
			buf = append(buf, node.Label...)
		}
		data[idx] = buf
		idx += 1
	}

	digest := new(big.Int).SetBytes(p.h.HashSlices(data))

	return digest, nil
}
