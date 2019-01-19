# rpost

Rational proofs of space time playgounrd and benchmarking

## Building

```
go get -u github.com/kardianos/govendor
govendor sync
go build
```

## Features
- [ ] Implement protocols store, prove and verify
  - [x] Store
  - [x] Prove
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

### Proof runtime profile
- 107.03s total running time for ReadProofs() - k=256
    - 103.09s in tree store read
        - 98.71s spent in file read syscall - ~95%
        - 4.38s in file indices computation - ~5%