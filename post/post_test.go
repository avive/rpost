package post

import (
	"crypto/rand"
	"github.com/avive/rpost/shared"
	"github.com/stretchr/testify/assert"
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

	n := uint64(6) // in bits. table size. T=2^n
	l := uint(20)   // in bits. difficulty. also # of nonce bits to store
	h := shared.NewHashFunc(id) // H(id) to be used for ipow

	table, err := NewTable(id, n, l, h, f)
	assert.NoError(t, err)

	// Create the file
	err = table.Generate()
	assert.NoError(t, err)
}
