package main

import (
	"math"
)

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
	var segmentSize int = end - start + 1
	var basePrimes []int = standardSieve(int(math.Sqrt(float64(end))))

	var isPrime []bool = make([]bool, segmentSize)

	for i := 0; i < len(basePrimes); i++ {
		p := basePrimes[i]

		st := ((start + p - 1) / p) * p // the closest bigger or equil number to `start`
		st -= start                     // gives the starting isPrime's index
		for j := st; j <= end; j += p {
			isPrime[j] = false
		}
	}

	for i := 0; i < segmentSize; i++ {
		if isPrime[i] {
			p := start + i
			for k := p * p; k <= end; k += p {
				isPrime[k] = false
			}
		}
	}

	var primes []int
	for i := 0; i <= segmentSize; i++ {
		if isPrime[i] {
			primes = append(primes, start+i)
		}
	}

	return primes
}
