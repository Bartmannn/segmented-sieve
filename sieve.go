package main

import (
	"math"
)

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

func standardSieve(end int) []int {
	isPrime := make([]bool, end+1)

	// TODO: Sprawdź czy można to zrobić lepiej,
	// poczytaj o różnych opcjach pętli for
	for i := 2; i <= end; i++ {
		isPrime[i] = true
	}

	for p := 2; p*p <= end; p++ {
		if isPrime[p] {
			for k := p * p; k <= end; k += p {
				isPrime[k] = false
			}
		}
	}

	var primes []int
	for i := 2; i <= end; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	return primes
}

func segmentedSieve(start, end int) []int {
	if start <= 2 {
		return standardSieve(end)
	}

	var segmentSize int = end - start + 1
	var basePrimes []int = standardSieve(int(math.Sqrt(float64(end))))

	var isPrime []bool = make([]bool, segmentSize)

	for i := 0; i < segmentSize; i++ {
		isPrime[i] = true
	}

	for _, p := range basePrimes {
		mul := (start + p - 1) / p
		st := mul * p // the closest bigger or equil number to `start`
		st -= start   // gives the starting isPrime's index

		// if mul is less than 2, then we should skip one tour
		if mul < 2 {
			st += p
		}

		for j := st; j < segmentSize; j += p {
			isPrime[j] = false
		}
	}

	for i := 0; i < segmentSize; i++ {
		if isPrime[i] {
			p := start + i
			for k := p * p; k < segmentSize; k += p {
				isPrime[k] = false
			}
		}
	}

	var primes []int
	for i := 0; i < segmentSize; i++ {
		if isPrime[i] {
			primes = append(primes, start+i)
		}
	}

	return primes
}

func genericSieve[T Number](end T) []T {
	isPrime := make([]bool, int(end)+1)

	// TODO: Sprawdź czy można to zrobić lepiej,
	// poczytaj o różnych opcjach pętli for
	for i := 2; i <= int(end); i++ {
		isPrime[i] = true
	}

	for p := T(2); p <= end/p; p++ {
		if isPrime[int(p)] {
			for k := p * p; k <= end; k += p {
				isPrime[k] = false
			}
		}
	}

	var primes []T
	for i := T(2); i <= end; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	return primes
}

func genericSegmentedSieve[T Number](start, end T) []T {
	if start <= 2 {
		return genericSieve(end)
	}

	var segmentSize T = end - start + 1
	var basePrimes []T = genericSieve(T(math.Sqrt(float64(end))))

	var isPrime []bool = make([]bool, segmentSize)

	for i := T(0); i < segmentSize; i++ {
		isPrime[i] = true
	}

	for _, p := range basePrimes {
		mul := (start + p - 1) / p
		st := mul * p // the closest bigger or equil number to `start`
		st -= start   // gives the starting isPrime's index

		// if mul is less than 2, then we should skip one tour
		if mul < 2 {
			st += p
		}

		for j := st; j < segmentSize; j += p {
			isPrime[j] = false
		}
	}

	for i := T(0); i < segmentSize; i++ {
		if isPrime[i] {
			p := start + i
			for k := p * p; k < segmentSize; k += p {
				isPrime[k] = false
			}
		}
	}

	var primes []T
	for i := T(0); i < segmentSize; i++ {
		if isPrime[i] {
			primes = append(primes, start+i)
		}
	}

	return primes
}
