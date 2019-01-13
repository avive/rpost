package post

import (
	"encoding/binary"
	"fmt"
	"github.com/avive/rpost/shared"
	"github.com/icza/bitio"
	"math"
	"math/big"
	"os"
)

const (
	fileBuffSizeBytes = 1024 * 1024 * 1024
	k                 = 256 // this can't be modified and is set as it is the bit length of the output of sha256()
)

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
	// w := internal.NewWriterSize(f, fileBuffSizeBytes)

	// bit writer around the file
	table.bw = bitio.NewWriter(f)

	return &table, nil
}

var bigOne = big.NewInt(1)

func (t *Table) Generate() error {

	n := uint64(math.Pow(2, float64(t.n)))
	fmt.Printf("Table size: %d \n", n)

	// compute probability in (0...1)
	p := GetProbability(t.l)
	fmt.Printf("Difficulty p: %.30f\n", p)

	fmt.Printf("Commitment x: 0x%x\n", t.id)
	fmt.Printf("Store file: %s\n", t.filePath)

	// TODO: bits to store is equals t.l - verify with tal
	// number of bites to store per hash
	// bits := uint(math.Ceil(math.Log2(1 / p)))
	// fmt.Printf("Nonce of bits to store : %d\n", bits)
	fmt.Printf("Number of nonce bits to store : %d\n", t.l)

	// create a bit mask of t.l bits set to 1
	storeMask := GetSimpleMask(t.l)
	fmt.Printf("Store mask bit field : %d %b\n", storeMask, storeMask)

	m := GetMask(32, t.l)
	fmt.Printf("Mask : %s\n", m.String())

	iBuf := make([]byte, 10)
	nonce := big.NewInt(0)
	d := new(big.Int)

	// todo: uint64 maxVal is 18446744073709551615 - is this ok or do we need to iterate over a 256bits big int here?
	for i := uint64(0); i < n; i++ {

		// nonce is in {0,1}^log(k/p) ????
		nonce = nonce.SetUint64(0)
		ln := binary.PutUvarint(iBuf, i)

		for {

			digest := t.h.Hash(iBuf[:ln], nonce.Bytes())
			d = d.SetBytes(digest)

			if d.Cmp(m) == -1 { // H(id, i, x) < p
				fmt.Printf(" Nonce: %d %b - digest: 0x%x\n", nonce.Uint64(), nonce.Uint64(), digest)

				// Pull exactly l lsb bits from nonce and store as uint64
				data := nonce.And(nonce, storeMask).Uint64()
				fmt.Printf("Data (%d lsb bits of nonce): %d %b \n", t.l, data, data)

				// Write the data to the file
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
