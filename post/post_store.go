package post

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/icza/bitio"
	"os"
)

// StoreWriter is a serial writer that appends data to the file
type StoreWriter interface {
	Write(r uint64, n byte) error
	WriteBool(b bool) error
	Close() error
}

// StoreReader is a random access reader capable of reading data from any valid bit offset
type StoreReader interface {
	Read(idx uint64) (bitarray.BitArray, error)
	ReadUint64(idx uint64) (uint64, error)
	ReadBytes(idx uint64) ([]byte, error)
	Close() error
}

type store struct {
	filePath string   // disk store data file path + name
	file     *os.File // file
	writer   bitio.Writer
	n        uint   // number of bits stored per entry
	sz       uint64 // file size in bytes - only used when reading
}

func NewStoreWriter(filePath string, n uint) (StoreWriter, error) {

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &store{filePath,
		f,
		bitio.NewWriter(f),
		n, 0}, nil
}

func NewStoreReader(filePath string, n uint) (StoreReader, error) {

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &store{filePath,
		f,
		nil,
		n,
		uint64(fi.Size())}, nil
}

func (s *store) Write(r uint64, n byte) error {
	return s.writer.WriteBits(r, n)
}

func (s *store) WriteBool(b bool) error {
	return s.writer.WriteBool(b)
}

func (s *store) Close() error {
	return s.writer.Close()
}

// Read from index id and return decoded uint64
func (s *store) ReadUint64(idx uint64) (uint64, error) {
	v, err := s.Read(idx)
	if err != nil {
		return 0, err
	}

	res, err := Uint64Value(v, uint64(s.n))
	if err != nil {
		return 0, err
	}

	return res, nil
}

// read from index idx and return as []byte
func (s *store) ReadBytes(idx uint64) ([]byte, error) {
	v, err := s.ReadUint64(idx)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, v)
	res := buf[:n]

	return res, nil
}

// Read s.n bits at index idx from the store
// This method can read from any index any number of times. e.g. It is not a classic GO reader which consumes the data it reads
// We need this for random access into the table when generating a proof for a random challenge...
func (s *store) Read(idx uint64) (bitarray.BitArray, error) {

	fmt.Printf("File size in bytes: %d...\n", s.sz)
	fmt.Printf("Reading %d bits entry from store at index %d...\n", s.n, idx)

	// First, figure out how many bytes we need to read and in which offset
	l := s.n / 8
	if s.n%8 != 0 {
		// we need to grow the buffer to accommodate up to 7 extra bits of data
		l += 1
	}

	offsetBits := idx * uint64(s.n)
	fmt.Printf("Bits offset: %d\n", offsetBits)

	offsetBytes := offsetBits / 8
	fmt.Printf("Bytes offset: %d\n", offsetBytes)

	if offsetBits%8 != 0 && offsetBytes+uint64(l)+1 < s.sz {
		// we are starting to read before the first data bit so we
		// need to grow the buffer by one unless this there are no more bytes in the file
		l += 1
	}

	// Read data goes here
	res := bitarray.NewBitArray(uint64(s.n), false)

	buff := make([]byte, l)
	n, err := s.file.ReadAt(buff, int64(offsetBytes))
	if err != nil {
		return res, err
	}
	if n == 0 {
		return res, errors.New("data for idx not found")
	}

	// Copy s.n bits of the data from the buffer at the correct positions

	// offset to read from current byte
	o := offsetBits % 8
	fmt.Printf("Initial read offset: %d\n", o)

	// current byte in buff to read from
	byteIdx := 0

	// create a bit array of s.n bits from the data in the buff
	// starting at position offsetBits % 8
	for i := uint64(0); i < uint64(s.n); i++ {

		// read next bit from the buffer
		// we need to use 7 - o because we assume byte representation of a bit field in the form: [7,6,5,...0]
		set := GetNthBit(buff[byteIdx], 7-o)
		if set {
			// set the bit in the ith index of the result bit array
			err := res.SetBit(i)
			if err != nil {
				return res, err
			}
		}
		o += 1
		if o == 8 { // moving forward...
			o = 0
			byteIdx += 1
		}
	}
	return res, nil
}
