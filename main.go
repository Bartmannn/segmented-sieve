package main

import (
	"bufio"
	"fmt"
	"os"
)

func saveToFile[T Number](filename string, primes []T) error {
	fmt.Println("Saving...")

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, p := range primes {
		fmt.Fprintf(writer, "%b %d\n", p, p)
	}

	fmt.Println("Completed!")

	return nil
}

func prepareRange[T Number](bits uint8) (T, T) {
	var start T = (1 << (bits - 1)) | 1            // 100..001
	var end T = (((1 << (bits - 1)) - 1) << 1) | 1 // 111..111
	return start, end
}

func runSieve(bits uint8) error {
	fmt.Println("Looking for prime numbers...")

	if bits <= 8 {
		start, end := prepareRange[uint8](bits)
		var primes []uint8 = genericSegmentedSieve(start, end)
		var filename string = fmt.Sprintf("primes_%db.txt", bits)
		return saveToFile(filename, primes)
	} else if bits <= 16 {
		start, end := prepareRange[uint16](bits)
		var primes []uint16 = genericSegmentedSieve(start, end)
		var filename string = fmt.Sprintf("primes_%db.txt", bits)
		return saveToFile(filename, primes)
	} else if bits <= 32 {
		start, end := prepareRange[uint32](bits)
		var primes []uint32 = genericSegmentedSieve(start, end)
		var filename string = fmt.Sprintf("primes_%db.txt", bits)
		return saveToFile(filename, primes)
	} else if bits > 64 {
		return fmt.Errorf("%d cannot be handled!", bits)
	}

	start, end := prepareRange[uint64](bits)
	var primes []uint64 = genericSegmentedSieve(start, end)
	var filename string = fmt.Sprintf("primes_%bd.txt", bits)
	return saveToFile(filename, primes)
}

func main() {
	var bits uint8 = 32
	runSieve(bits)
}
