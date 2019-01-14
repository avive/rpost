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
	/*
	testPost(t, 4, 20)
	testPost(t, 4, 22)
	testPost(t, 8, 18)
	testPost(t, 10, 16)*/

	testPost(t, 12, 20)
}

// n - Table size T = 2^n
// l - iPoW difficulty and the # of nonce bits to store per entry
func testPost(t *testing.T, n uint64, l uint) {

	// File to store iPoWs
	f := filepath.Join(os.TempDir(), "post.bin")

	// Initial commitment
	id := make([]byte, 32)
	_, err := rand.Read(id)
	assert.NoError(t, err)

	h := shared.NewHashFunc(id) // H(id) to be used for iPoW

	table, err := NewTable(id, n, l, h, f)
	assert.NoError(t, err)

	// Create the file
	err, res := table.Generate(true)
	assert.NoError(t, err)

	err = dumpContent(f, l)
	assert.NoError(t, err)

	storeReader, err := NewStoreReader(f, l)
	assert.NoError(t, err)

	tableSize := uint64(math.Pow(2, float64(n)))
	validateStoreSize(t, f, tableSize, uint64(l))

	for i := uint64(0); i < tableSize; i++ {
		data, err := storeReader.Read(i)
		assert.NoError(t, err, "index: %d", i)

		// compare the data parsed from the file to the data
		// returned by Generate in
		v, err := Uint64Value(data, uint64(l))
		assert.NoError(t, err)
		assert.Equal(t, v, res[i])

		s, err := String(data, uint64(l))
		fmt.Printf("Data: %s \n", s)
	}
}

// Validate actual store file size based on expected values
func validateStoreSize(t *testing.T, filePath string, tableSize uint64, bitsPerEntry uint64) {
	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()
	fileInfo, err := file.Stat()
	assert.NoError(t, err)
	expectedFileSize := tableSize*bitsPerEntry/8 + (tableSize % 8)
	assert.Equal(t, expectedFileSize, uint64(fileInfo.Size()))
}
