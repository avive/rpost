## Proof random file reads analysis

#### Definitions
- k := security param. e.g. 128
- n := log2 of # of table entries. e.g. T = 2^n 
- n is also height of merkle proof of store data (incl entrie leaf sibling)
- p* := prob to find a path probe in a try = k / T = k / 2^n (paper denotation)
- e := expected # of tries to find a valid path probe = 1 / p* = 2^n / k (paper denotation). See page 7: Ω(T/k) attempts


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
          
#### An example
Let n=20 and k=128.
- Table has 1,048,576 data entries 
- 8192 expected ops to find a valid path probe
- Computing `2^n * n * k` we get: 2^20 * 20 * 128 = `2,684,354,560 random access i/o ops`


