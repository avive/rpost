# rpost

Rational proofs of space time playgounrd and benchmarking

## Building

```
go get -u github.com/kardianos/govendor
govendor sync
go build
```

## Features
- [ ] Implement the protocol - store, prove and verify
  - [x] Store
  - [ ] Prove
  - [ ] Verify
- [x] Optimal store size
- [x] Support all paper params
- [x] Fast random-access of data from store (bit-level)
- [x] Table generation and validity tests
- [x] Tests using in-memory table data
- [x] Optimal Merkle tree generation and store 
- [ ] Real-world test scenarios

## Testing
```
go test ./...
```

## Profiling notes
- 107.03s total ReadPRoofs()
- 101.28s in file read syscall
- 95% of time in file read
- Total running time gen proof: 2m32s... for n=10, l=10
