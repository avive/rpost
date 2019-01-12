package post

import (
	"encoding/binary"
	"fmt"
	"github.com/avive/rpost/internal"
	"github.com/avive/rpost/shared"
	"github.com/icza/bitio"
	"math"
	"math/big"
	"os"
)

const fileBuffSizeBytes = 1024 * 1024 * 1024

type Table struct {
	id       []byte          // initial commitment
	n        uint64          // n param 1 <= n <= 63 - table size is 2^n
	l        uint            // l param (num of leading 0s for p) := f(p). 1: 50%, 2: 25%, 3:12.5%...
	h        shared.HashFunc // Hx()
	filePath string          // disk store data file path + name
	file     *os.File        // file
	bw       bitio.Writer    // Bitio around file writer
}

// Create a new prover with commitment X and param 1 <= n <= 63
func NewTable(id []byte, n uint64, l uint, h shared.HashFunc, filePath string) (*Table, error) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	table := Table{id, n, l, h, filePath, f, nil}

	// buffered io around a file
	w := internal.NewWriterSize(f, fileBuffSizeBytes)

	// bit writer around the buffer
	table.bw = bitio.NewWriter(w)

	return &table, nil
}

var bigOne = big.NewInt(1)

func (t *Table) Generate() error {

	n := uint64(math.Pow(2, float64(t.n)))
	fmt.Printf("Table size: %d \n", n)

	// compute probability in (0...1)
	p := difficultyToProb(t.l)
	fmt.Printf("Difficulty p: %.30f\n", p)

	fmt.Printf("Commitment x: 0x%x\n", t.id)
	fmt.Printf("Store %s\n", t.filePath)

	// number of bites to store per hash
	// bits := uint(math.Ceil(math.Log2(1 / p)))
	// fmt.Printf("Nonce of bits to store : %d\n", bits)

	// bits to store is equals t.l !!!
	fmt.Printf("Number of nonce bits to store : %d\n", t.l)

	// Generate a mask here of 32 bytes with l 0 bits at msb
	mask := make([]byte, 32)
	for i := 0; i < 32; i++ {
		mask[i] = 0xff
	}

	m := clearMSBBits(t.l, mask)
	mask = m.Bytes()

	iBuf := make([]byte, 10)
	nonce := big.NewInt(0)
	d := new(big.Int)

	// maxVal is 18446744073709551615
	for i := uint64(0); i < n; i++ {

		// nonce is in {0,1}^log(k/p) ????
		nonce = nonce.SetUint64(0)
		ln := binary.PutUvarint(iBuf, i)

		for {

			digest := t.h.Hash(iBuf[:ln], nonce.Bytes())
			d = d.SetBytes(digest)

			if d.Cmp(m) == -1 { // H(id, i, x) < p
				fmt.Printf(" Nonce: %d - digest: 0x%x\n", nonce.Uint64(), digest)

				// todo: we want to set data so it is t.l LSB bits of nonce...
				// nonce => data

				// big endian uint64 from big int
				data := nonce.Uint64()

				err := t.bw.WriteBits(data, byte(t.l))
				if err != nil {
					return err
				}
				break
			}

			nonce = nonce.Add(nonce, bigOne)
		}
	}

	return t.finalize()
}

func (t *Table) finalize() error {
	return t.bw.Close()
}
