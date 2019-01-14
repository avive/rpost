package post

import (
	"crypto/rand"
	"fmt"
	"github.com/avive/rpost/shared"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestPost(t *testing.T) {

	// File to store ipows
	f := filepath.Join(os.TempDir(), "post.bin")

	// Initial commitment
	id := make([]byte, 32)
	_, err := rand.Read(id)
	assert.NoError(t, err)

	n := uint64(2)              // in bits. table size. T=2^n
	l := uint(20)               // in bits. difficulty. also # of nonce bits to store
	h := shared.NewHashFunc(id) // H(id) to be used for iPoW

	table, err := NewTable(id, n, l, h, f)
	assert.NoError(t, err)

	// Create the file
	err = table.Generate()
	assert.NoError(t, err)

	err = dumpContent(f, l)
	assert.NoError(t, err)

	storeReader, err := NewStoreReader(f, l)
	assert.NoError(t, err)

	tableSize := uint64(math.Pow(2, float64(n)))
	for i := uint64(0); i < tableSize; i++ {
		data, err := storeReader.Read(i)
		assert.NoError(t, err, "index: %d", i)
		s, err := String(data, uint64(l))
		fmt.Printf("Data: %s \n", s)
	}
}
