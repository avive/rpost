package merkle

import (
	"errors"
	"math"
	"os"
)

// A simple (k,v) store for storing indexed labels

// todo: review these
const (
	W             = 256         // label length in bits
	WB            = 32          // W length in bytes
	buffSizeBytes = 1024 * 1024 // Write buffer size
)

type Label []byte      // label is WB bytes long binary data
type Labels []Label    // an ordered list of Labels
type Identifier string // A Variable-length binary string. e.g. "0011010" Only 0s and 1s are allowed chars.
type WriteData struct {
	id Identifier
	l  Label
}

// A simple (k,v) store writer
// Labels must be written in depth-first order. Random access is not supported
type IKvStoreWriter interface {
	Write(id Identifier, l Label)
	IsLabelInStore(id Identifier) (bool, error)
	Reset() error
	Delete() error
	Size() uint64
	Finalize()    // finalize writing w/o closing the file
	Close() error // finalize and close
}

// A simple (k,v) reader - fully supports random access
type IKvStoreReader interface {
	Read(id Identifier) (Label, error)
	Size() uint64
	Close() error
}

type kvFileStore struct {
	fileName string
	file     *os.File
	n        uint // 1 <= n < 64
	f        BinaryStringFactory
	bw       *Writer
	c        uint64 // num of labels written to store in this session
}

// n specifies the leafs height from the root
func NewKvFileStoreWriter(fileName string, n uint) (IKvStoreWriter, error) {
	res := &kvFileStore{
		fileName: fileName,
		n:        n,
		f:        NewSMBinaryStringFactory(),
	}

	f, err := os.OpenFile(res.fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	res.file = f
	res.bw = NewWriterSize(f, buffSizeBytes)
	return res, err
}

func NewKvFileStoreReader(fileName string, n uint) (IKvStoreReader, error) {
	res := &kvFileStore{
		fileName: fileName,
		n:        n,
		f:        NewSMBinaryStringFactory(),
	}

	f, err := os.OpenFile(res.fileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	res.file = f
	return res, err
}

func (d *kvFileStore) Write(id Identifier, l Label) {
	d.c += 1
	_, err := d.bw.Write(l)
	if err != nil {
		panic(err)
	}
}

// Removes all data from the file
func (d *kvFileStore) Reset() error {
	err := d.bw.Flush()
	if err != nil {
		return err
	}

	d.c = 0
	return d.file.Truncate(0)
}

func (d *kvFileStore) Finalize() {
	// flush buffer to file
	if d.bw != nil {
		_ = d.bw.Flush()
	}
}

func (d *kvFileStore) Close() error {
	d.Finalize()
	return d.file.Close()
}

func (d *kvFileStore) Delete() error {
	return os.Remove(d.fileName)
}

func (d *kvFileStore) Size() uint64 {
	stats, err := d.file.Stat()
	if err != nil {
		println(err)
	}

	res := uint64(stats.Size())

	if d.bw != nil {
		res += uint64(d.bw.Buffered())
	}

	return res
}

// Returns true iff node's label is already the store
func (d *kvFileStore) IsLabelInStore(id Identifier) (bool, error) {

	idx, err := d.calcFileIndex(id)
	if err != nil {
		return false, err
	}

	stats, err := d.file.Stat()
	if err != nil {
		return false, err
	}

	if d.bw.Buffered() > 0 && idx < (d.c*WB) {
		// label is in file or in the buffer
		return true, nil
	}

	fileSize := uint64(stats.Size())
	return idx < fileSize, nil
}

// Read label value from the store
// Returns the label of node id or error if it is not in the store
func (d *kvFileStore) Read(id Identifier) (Label, error) {

	label := make(Label, WB)

	// total # of labels written - # of buffered labels == idx of label at buff start
	// say 4 labels were written, and Buffered() is 64 bytes. 2 last labels
	// are in buffer and the index of the label at buff start is 2.
	// idAtBuffStart := d.c - uint64(d.bw.Buffered()/shared.WB)

	// label file index
	idx, err := d.calcFileIndex(id)
	if err != nil {
		return label, err
	}

	n, err := d.file.ReadAt(label, int64(idx))
	if err != nil {
		return label, err
	}

	if n == 0 {
		return label, errors.New("label for id is not in store")
	}

	return label, nil
}

// Returns the file offset for a node id
func (d *kvFileStore) calcFileIndex(id Identifier) (uint64, error) {
	s := d.subtreeSize(id)
	s1, err := d.leftSiblingsSubtreeSize(id)
	if err != nil {
		return 0, err
	}

	idx := s + s1 - 1
	offset := idx * WB
	//fmt.Printf("Node id %s. Index: %d. Offset: %d\n", id, idx, offset)
	return offset, nil
}

// Returns the size of the subtree rooted at node id
func (d *kvFileStore) subtreeSize(id Identifier) uint64 {
	// node depth is the number of bits in its id
	depth := uint(len(id))
	height := d.n - depth
	return uint64(math.Pow(2, float64(height+1)) - 1)
}

// Returns the size of the subtrees rooted at left siblings on the path
// from node id to the root node
func (d *kvFileStore) leftSiblingsSubtreeSize(id Identifier) (uint64, error) {
	bs, err := d.f.NewBinaryString(string(id))
	if err != nil {
		return 0, err
	}

	siblings, err := bs.GetBNSiblings(true)
	if err != nil {
		return 0, err
	}
	var res uint64

	for _, s := range siblings {
		res += d.subtreeSize(Identifier(s.GetStringValue()))
	}

	return res, nil
}
