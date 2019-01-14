package post

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/avive/rpost/shared"
	"math"
	"math/big"
	"math/bits"
)

const (
	k = 256 // this can't be modified and is set as it is the bit length of the output of sha256()
)

type Table struct {
	id []byte          // initial commitment
	n  uint64          // n param 1 <= n <= 63 - table size is 2^n
	l  uint            // l param (num of leading 0s for p) := f(p). 1: 50%, 2: 25%, 3:12.5%...
	h  shared.HashFunc // Hx()
	s  StoreWriter
}

// Create a new prover with commitment X and param
// n:= 1 <= n <= 63
// l:= 1 <= l <= 63
func NewTable(id []byte, n uint64, l uint, h shared.HashFunc, filePath string) (*Table, error) {

	// todo: check precondition args values here

	fmt.Printf("Store file: %s\n", filePath)

	store, err := NewStoreWriter(filePath, l)
	if err != nil {
		return nil, err

	}
	table := Table{id, n, l, h, store}
	return &table, nil
}

var one = big.NewInt(1)

// var maxNonce = GetMaxNonce(256)

// Implements the Store phase of rpost (page 9)
func (t *Table) Store() error {

	// 1. Generate and store the values of the iPoW table G
	err, _ := t.Generate(false)
	return err

	// todo: 2. Compute commitment com on G (root of Merkle tree where G data are leaves
}

func (t *Table) Generate(returnData bool) (error, []uint64) {

	n := uint64(math.Pow(2, float64(t.n)))
	fmt.Printf("Table size: %d \n", n)

	// p*
	phi := k / float64(n)
	fmt.Printf("P*: %f \n", phi)

	// compute probability in (0...1)
	p := GetProbability(t.l)
	fmt.Printf("Difficulty p: %.30f\n", p)

	fmt.Printf("Expected hashes to find a digest is at least : %d \n", int(1/p))

	maxNonceVal := big.NewInt(int64(math.Ceil(k / p)))
	fmt.Printf("Max permitted nonce: %s\n", maxNonceVal.String())

	fmt.Printf("Commitment x: 0x%x\n", t.id)

	// number of bites to store per hash is ame as l
	//bits := uint(math.Ceil(math.Log2(1 / p))) === t.l
	fmt.Printf("Number of nonce bits to store : %d\n", t.l)
	fmt.Printf("Difficulty param : %d\n", t.l)

	// create a bit mask of t.l bits set to 1
	storeMask := GetSimpleMask(t.l)
	fmt.Printf("Store mask bit field : %d %b\n", storeMask, storeMask)

	m := GetMask(32, t.l)
	fmt.Printf("Mask : %s\n", m.String())

	iBuf := make([]byte, 10)
	nonce := big.NewInt(0)
	d := new(big.Int)

	var res []uint64

	for i := uint64(0); i < n; i++ {

		// nonce is in {0,1}^log(k/p) - max nonce value is k/p
		nonce = nonce.SetUint64(0)
		ln := binary.PutUvarint(iBuf, i)

		for {

			digest := t.h.Hash(iBuf[:ln], nonce.Bytes())
			d = d.SetBytes(digest)

			if d.Cmp(m) == -1 { // H(id, i, x) < p
				fmt.Printf(" Nonce: %d %b - digest: 0x%x\n", nonce.Uint64(), nonce.Uint64(), digest)

				// Take l lsb bits from nonce and store as uint64
				data := nonce.And(nonce, storeMask).Uint64()

				fmt.Printf("Data (%d lsb bits of nonce): %d %b bits:%d \n", t.l, data, data, bits.Len64(data))

				// Write the data to the file - exactly t.l lsb bits of data
				// if t.l > len(data) then 0s are padded starting MSB bit
				// so, for example, if len(data) = 16 and t.l = 20, 4 leading 0s will be written starting at MSB bit (left-to-right)
				// and the 16 bits of data next using big-endian encoding. e.g. MSB bit first...
				err := t.s.Write(data, byte(t.l))
				if err != nil {
					return err, nil
				}

				if returnData {
					res = append(res, data)
				}

				break
			}

			nonce = nonce.Add(nonce, one)

			if nonce.Cmp(maxNonceVal) == 1 {
				// nonce overflow case. We expect nonce to be up to ceil(k/p)
				return errors.New("failed to find nonce in permitted range"), nil
			}
		}
	}

	return t.finalize(), res
}

func (t *Table) finalize() error {
	return t.s.Close()
}
