## Proof random file reads analysis

#### Definitions
- k := 256 (constant)
- n := log2 of # of table entries. e.g. T = 2^n 
- n is also height of merkle proof of store data (incl entrie leaf sibling)
- p* := prob to find a path probe in a try = k / T = k / 2^n (paper denotation)
- e := expected # of tries to find a valid path probe = 1 / p* = 2^n / k (paper denotation)


#### Computing the expected number of random file access read ops

From the proof protocol we get: 

- k * e * k * n = 
- k^2 * e * n = 
- (2^n / k) * n * k^2 = 
- n * 2^n * k

Reasoning:
- In each of the `k iterations` we perform an expected `e iterations`
    - In each `e iteration` we`read k values` from the data store 
      - for reach `k value` we read `n values` to compute a merkle proof
          
#### A concrete example
- Let n=20 and k=256
- Table has 1,048,576 data entries 
- 64 expected ops to find a valid path probe
- Computing `2^n * n * k` we get: 2^14 * 14 * 256 = `5,368,709,120 expected random access i/o ops`
