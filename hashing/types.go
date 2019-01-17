package hashing

type HashFunc interface {
	// Hash takes arbitrary binary data and returns WB bytes
	Hash(data ...[]byte) []byte
	HashSlices(data [][]byte) []byte
	HashSingle(data []byte) []byte
}
