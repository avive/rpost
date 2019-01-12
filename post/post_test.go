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

	f := filepath.Join(os.TempDir(), "post.bin")

	// set commitment to 32 random bytes
	id := make([]byte, 32)
	_, err := rand.Read(id)
	assert.NoError(t, err)

	n := uint64(10) // in bits
	l := uint(20)   // in bits
	h := shared.NewHashFunc(id)

	table, err := NewTable(id, n, l, h, f)
	assert.NoError(t, err)

	err = table.Generate()
	assert.NoError(t, err)
}
