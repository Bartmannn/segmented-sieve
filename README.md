# Segmented sieve

An educational implementation of the segmented Sieve of Eratosthenes in Go. The project searches for prime numbers with a given bit length, uses a `2-3-5` wheel to skip obvious composite candidates, and writes results in the following format:

```text
<number in binary> <number in decimal>
```

This is not a finished CLI tool yet. I treat it as a learning project for refreshing Go syntax, practicing generics, working with data structures and goroutines, and thinking more deliberately about performance.

## What The Program Does

For a given number of bits, the program builds the range:

```text
start = 100...001
end   = 111...111
```

For example, for 8-bit numbers it searches the range from `129` to `255`, so the values actually require 8 bits and the range boundaries are odd.

Then the code:

- generates base primes up to `sqrt(end)`,
- splits the searched range into segments,
- processes each segment with the generic `genericSegmentedSieve[T Number]` function,
- skips multiples of `2`, `3`, and `5` using a `2-3-5` wheel,
- can distribute segments across workers,
- writes the result to a file through a buffered writer.

## The 2-3-5 Wheel

The `2-3-5` wheel uses the fact that any prime candidate greater than `5` cannot be divisible by `2`, `3`, or `5`. Within modulo `30`, only these residues remain:

```text
1, 7, 11, 13, 17, 19, 23, 29
```

Instead of checking every next integer, the implementation walks through candidates using these steps:

```text
6, 4, 2, 4, 2, 4, 6, 2
```

This reduces the number of iterations and the amount of composite marking work. In the code, this is represented by the `Wheel235` structure, which stores the wheel steps and maps residues modulo `30` to wheel indices.

## Concurrency And Memory

The project experiments with processing segments through goroutines. A producer creates `Job[T]` values, workers run the sieve on their assigned ranges, and results are sent back through a channel as `Result` values.

The segment size is selected by default from the L1 data cache size reported by `github.com/klauspost/cpuid`. This is intentionally practical: the segment should be small enough to fit well in fast cache memory and reduce the cost of cache misses.

File output is buffered. For larger ranges this matters, because I/O can start to dominate the cost of the computation itself.

## Generics

The sieve is written generically for integer types:

```go
type Number interface {
    int | int8 | int16 | int32 | int64 |
        uint | uint8 | uint16 | uint32 | uint64
}
```

Because of that, the same core algorithm can work with `uint8`, `uint16`, `uint32`, `uint64`, and other integer types. This was one of the main reasons for building the project: moving from a simple sieve toward a version that uses Go generics in practice.

## Example Output For 8 Bits

Example prime numbers from the 8-bit range:

```text
10000011 131
10001001 137
10001011 139
10010101 149
10010111 151
10011101 157
10100011 163
10100111 167
10101101 173
10110011 179
10110101 181
10111111 191
11000001 193
11000101 197
11000111 199
11010011 211
11011111 223
11100011 227
11100101 229
11101001 233
11101111 239
11110001 241
11111011 251
```

## Tests

The project includes tests that compare the segmented implementation against expected results and a simpler naive oracle:

```bash
go test ./...
```

The current test suite passes for the existing code.

## Status And Future Work

This is still a work in progress. The most important planned improvements are:

- replacing the `[]bool` array with an integer-backed bitset, where each bit represents a `true/false` state,
- reducing stored states to only the candidates that are not skipped by the `2-3-5` wheel,
- adding structured benchmarks comparing the standard, segmented, generic, wheel-based, bitset-based, and different segment-size versions,
- measuring the impact of worker count, output buffer size, and cache-aware segment sizing,
- improving the execution interface so parameters do not have to be edited manually in `main`.

## Why This Project Exists

This code is not meant to pretend to be a cryptographic library or a production-ready tool. Its purpose is to show a learning process: from the classic Sieve of Eratosthenes, through a segmented sieve, toward optimizations involving wheels, cache behavior, generics, and concurrency.

The most interesting part of the project is not just finding prime numbers. It is about checking how specific implementation decisions affect performance and memory usage.
