package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// fmt.Println(standardSieve(34787))
	// fmt.Println(genericSegmentedSieve[uint32](2147483649, 4294967295))
	// fmt.Println(genericSegmentedSieve[int8](5, 127))
	saveToFile("primes.txt", genericSegmentedSieve[uint32](2147483649, 4294967295))
}

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

	fmt.Println("Primes numbers are saved!")

	return nil
}
