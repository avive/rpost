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

func TestPostGeneration(t *testing.T) {

	generatePost(t, 4, 20, "post1.bin", "merkle1.bin")

}


// n - Table size T = 2^n
// l - iPoW difficulty and the # of nonce bits to store per entry
// This is a playground disguised as a test :-)
func generatePost(t *testing.T, n uint64, l uint, fileName string, merkleFileName string) {

	// File to store iPoWs
	currFolder, err := os.Getwd() // os.TempDir()
	if err != nil {
		assert.NoError(t, err, "can't get path of executable")
	}

	f := filepath.Join(currFolder, fileName)
	mf := filepath.Join(currFolder, merkleFileName)

	// Initial commitment
	id := util.Rnd(t, 32)

	// H(id) to be used for iPoW
	h := hashing.NewHashFunc(id)

	// New store table
	table, err := NewTable(id, n, l, h, f)
	assert.NoError(t, err)

	// Store the data
	comm, err := table.Store(mf)
	assert.NoError(t, err)

	fmt.Printf("Merkle commitment: 0x%x \n", comm)

}
