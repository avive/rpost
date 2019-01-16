package post

import (
	"fmt"
	"github.com/avive/rpost/hashing"
	"github.com/avive/rpost/util"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestMerkleWriter(t *testing.T) {
	testMerkleStore(t, 4, 20, "post1.bin", "merkle1.bin")
}

// n - Table size T = 2^n
// l - iPoW difficulty and the # of nonce bits to store per entry
// This is a playground disguised as a test :-)
func testMerkleStore(t *testing.T, n uint64, l uint, postFileName string, merkleFileName string) {

	// File to store iPoWs
	currFolder, err := os.Getwd() // os.TempDir()
	if err != nil {
		assert.NoError(t, err, "can't get path of executable")
	}

	f := filepath.Join(currFolder, postFileName)
	mf := filepath.Join(currFolder, merkleFileName)

	// Initial commitment
	id := util.Rnd(t, 32)

	// H(id) to be used for iPoW
	h := hashing.NewHashFunc(id)

	// New store table
	tbl, err := NewTable(id, n, l, h, f)
	assert.NoError(t, err)

	// Store the data
	res, err := tbl.Generate(true)
	assert.NoError(t, err)

	// post memory reader from post data in ram
	sr := NewMemoryStoreReader(res)

	mw, err := NewMerkleTreeWriter(sr, mf, l, uint(n), h)
	assert.NoError(t, err)

	comm, err := mw.Write()
	assert.NoError(t, err)
	fmt.Printf("Merkle commitment: 0x%x \n", comm)

	// test merkle tree from post store
	sr, err = NewStoreReader(f, l)
	assert.NoError(t, err)

	mw, err = NewMerkleTreeWriter(sr, mf, l, uint(n), h)
	assert.NoError(t, err)
	comm1, err := mw.Write()
	assert.NoError(t, err)
	fmt.Printf("Merkle commitment: 0x%x \n", comm)

	assert.EqualValues(t, comm, comm1)

	mr, err := NewMerkleTreeReader(mf, l, uint(n-1), h)
	assert.NoError(t, err)

	path, err := mr.ReadPath("101")
	assert.NoError(t, err)
	for _, n := range path {
		fmt.Printf("Id: %s. Label: %x\n", n.id, n.label)
	}

	// close the reader when we are done with it
	err = mr.Close()
	assert.NoError(t, err)

}
