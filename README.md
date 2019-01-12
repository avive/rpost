# rpost

Rational proofs of space time playgounrd and benchmarking

## Building

```
go get -u github.com/kardianos/govendor
govendor sync
go build
```

## Testing
```
go test ./...
```

```
name  def	    Formula	    example     value	Range	Notes
k	Security param 256		
H(x)	Hash funciton {0,1}^*=>{0,1}^k	sha256(x)		{0,1}^k	
t	Log 2 of # of entries		20	>1	
T	Table size in entries (from t)	2^t	2^20		
p	Difficulty param		0.00001	(0,1)	
l	Number of leading 0s in {0,1}^k for p		5	(0...k)	When h={0,1}^k is considered 0.h
					
s	POST storage size bits	ceil(log2(1/p)		(0...k)	
TS	Total storage leafs	T*s == 2^t * s			
TSM	Total storage  merkle nodes	(2^t -1) * k			
TST	Total storage	2^t*s + (2^t -1)*k
```
