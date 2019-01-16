package post

import (
	"fmt"
	"github.com/avive/rpost/hashing"
	"github.com/avive/rpost/util"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestPost(t *testing.T) {

	// testPostStore(t, 4, 20, "post.bin")

	testPost(t, 4, 20, "post1.bin", "merkle1.bin")
	// testPost(t, 4, 22)
	// testPost(t, 8, 18)
	// testPost(t, 10, 16)

	// testPost(t, 12, 20)
}

// n - Table size T = 2^n
// l - iPoW difficulty and the # of nonce bits to store per entry
// This is a playground disguised as a test :-)
func testPost(t *testing.T, n uint64, l uint, fileName string, merkleFileName string) {

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


// n - Table size T = 2^n
// l - iPoW difficulty and the # of nonce bits to store per entry
// This is a playground disguised as a test :-)
func testPostStore(t *testing.T, n uint64, l uint, fileName string) {

	// File to store iPoWs
	currFolder, err := os.Getwd() // os.TempDir()
	if err != nil {
		assert.NoError(t, err, "can't get path of executable")
	}

	f := filepath.Join(currFolder, fileName)

	// Initial commitment
	id := util.Rnd(t, 32)

	// H(id) to be used for iPoW
	h := hashing.NewHashFunc(id)

	// New store table
	table, err := NewTable(id, n, l, h, f)
	assert.NoError(t, err)

	// Store the data
	res, err := table.Generate(true)
	assert.NoError(t, err)

	// Display the stored content from file
	err = dumpContent(f, l)
	assert.NoError(t, err)

	// test reading stored data from disc vs. expected data
	// returned in ram from Store()
	storeReader, err := NewStoreReader(f, l)
	assert.NoError(t, err)

	tableSize := uint64(math.Pow(2, float64(n)))
	validateStoreSize(t, f, tableSize, uint64(l))

	for i := uint64(0); i < tableSize; i++ {
		data, err := storeReader.Read(i)
		assert.NoError(t, err, "index: %d", i)

		// compare the data parsed from the file to the data
		// returned by Generate in
		v, err := util.Uint64Value(data, uint64(l))
		assert.NoError(t, err)
		assert.Equal(t, v, res[i])

		s, err := util.String(data, uint64(l))
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
