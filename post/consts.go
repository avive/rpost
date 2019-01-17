package post

// protocol hard-coded shared constants

const (
	K  = 256 // this needs to match the output size of Hx() - when sha256() is used
	W  = K   // merkle tree label length in bits
	WB = 32  // W length in bytes
)
