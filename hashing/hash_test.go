package hashing

import (
	"bytes"
	"fmt"
	"github.com/spacemeshos/sha256-simd"
	"math"
	"testing"
	"time"
)

func BenchmarkSha256(t *testing.B) {
	buff := bytes.Buffer{}
	buff.Write([]byte("Seed data goes here"))
	out := [32]byte{}
	n := uint64(math.Pow(10, 8))

	fmt.Printf("Computing %d serial sha-256s...\n", n)

	t1 := time.Now()

	for i := uint64(0); i < n; i++ {
		out = sha256.Sum256(buff.Bytes())
		buff.Reset()
		buff.Write(out[:])
	}

	e := time.Since(t1)
	r := n / (uint64(e.Seconds()))
	fmt.Printf("Final hash: %x. Running time: %s secs. Hash-rate: %d hashes-per-sec\n", buff.Bytes(), e, r)
}
