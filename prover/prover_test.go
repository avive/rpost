package prover

import (
	"fmt"
	"github.com/avive/rpost/hashing"
	"github.com/avive/rpost/post"
	"github.com/avive/rpost/util"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProver(t *testing.T) {
	testProver(t, 11, 14, "post1.bin", "merkle1.bin")
}

// n - Table size T = 2^n
// l - iPoW difficulty and the # of nonce bits to store per entry
func testProver(t *testing.T, n uint64, l uint, postFileName string, merkleFileName string) {

	currFolder, err := os.Getwd()
	if err != nil {
		assert.NoError(t, err, "can't get path of executable")
	}

	f := filepath.Join(currFolder, postFileName)
	mf := filepath.Join(currFolder, merkleFileName)

	// Initial commitment
	//id := util.Rnd(t, 32)
	seed, _ := new(big.Int).SetString("3b05a45e418666973c19aaccdf2547ba8d33e9610f547b31a0735d95d45469b5", 16)
	id := seed.Bytes()

	// H(id) to be used for iPoW
	h := hashing.NewHashFunc(id)

	// Generate a new store table

	tbl, err := post.NewTable(id, n, l, h, f)
	assert.NoError(t, err)
	_, err = tbl.Generate(false)
	assert.NoError(t, err)

	// Generate merkle tree from post store
	sr, err := post.NewStoreReader(f, l)
	assert.NoError(t, err)
	mw, err := post.NewMerkleTreeWriter(sr, mf, l, uint(n), h)
	assert.NoError(t, err)
	comm, err := mw.Write()
	assert.NoError(t, err)

	fmt.Printf("Merkle commitment: 0x%x \n", comm)


	// Generate a proof for a challenge

	pv, err := NewProver(id, n, l, h, f, mf)

	challenge := util.Rnd(t, 32)

	t1 := time.Now()
	proof, err := pv.Prove(challenge)
	e1 := time.Since(t1)
	t.Logf("Proof generated in %s seconds.\n", e1)

	assert.NoError(t, err)
	for i, n := range proof.Nonces {
		fmt.Printf("[%d] : %d\n", i, n)
	}

	assert.NoError(t, err)
}
