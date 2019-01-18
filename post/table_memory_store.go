package post

import (
	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/avive/rpost/util"
)

// a simple in-ram post data store implementing StoreReader
type MemoryStore struct {
	data []uint64
}

func (ms *MemoryStore) ReadUint64(idx uint64) (uint64, error) {
	return ms.data[idx], nil
}

func NewMemoryStoreReader(data []uint64) StoreReader {
	return &MemoryStore{data}
}

func (ms *MemoryStore) Read(idx uint64) (bitarray.BitArray, error) {
	panic("not yet implemented")
}

func (ms *MemoryStore) Close() error {
	return nil
}

func (ms *MemoryStore) FileName() string {
	return "" // hack alert
}

// read from index idx and return as []byte
func (ms *MemoryStore) ReadBytes(idx uint64) ([]byte, error) {
	v, err := ms.ReadUint64(idx)
	if err != nil {
		return nil, err
	}

	return util.EncodeToBytes(v), nil

}
