package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/klauspost/cpuid"
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

func runSieve(bits, workersNum uint8, segmentSize int) error {
	fmt.Println("Looking for prime numbers...")
	return nil
	// if bits <= 8 {
	// 	start, end := prepareRange[uint8](bits)
	// 	// jobs := make(chan Job[uint8])
	// 	var primes []uint8 = genericSegmentedSieve(start, end)
	// 	var filename string = fmt.Sprintf("primes_%db.txt", bits)
	// 	return saveToFile(filename, primes)
	// } else if bits <= 16 {
	// 	start, end := prepareRange[uint16](bits)
	// 	// jobs := make(chan Job[uint16])
	// 	var primes []uint16 = genericSegmentedSieve(start, end)
	// 	var filename string = fmt.Sprintf("primes_%db.txt", bits)
	// 	return saveToFile(filename, primes)
	// } else if bits <= 32 {
	// 	start, end := prepareRange[uint32](bits)
	// 	// jobs := make(chan Job[uint32])
	// 	var primes []uint32 = genericSegmentedSieve(start, end)
	// 	var filename string = fmt.Sprintf("primes_%db.txt", bits)
	// 	return saveToFile(filename, primes)
	// } else if bits > 64 {
	// 	return fmt.Errorf("%d cannot be handled!", bits)
	// }

	// start, end := prepareRange[uint64](bits)
	// // jobs := make(chan Job[uint64])
	// var primes []uint64 = genericSegmentedSieve(start, end)
	// var filename string = fmt.Sprintf("primes_%bd.txt", bits)
	// return saveToFile(filename, primes)
}

func run(bits, workersNum uint8, segmentSize int, filename string, bufferSizeMB int) error {
	results := make(chan Result, int(workersNum)*2)
	var newFilename string = filename
	var wg sync.WaitGroup

	if bufferSizeMB == 0 {
		bufferSizeMB = 4
	}

	if filename == "" {
		newFilename = fmt.Sprintf("primes_%db.txt", bits)
	}

	if bits <= 8 {
		start, end := prepareRange[uint8](bits)
		var primes []uint8 = genericSegmentedSieve(start, end)
		fmt.Println(primes)
		return nil
	} else if bits <= 16 {
		start, end := prepareRange[uint16](bits)
		jobs := make(chan Job[uint16])

		go func() {
			defer close(jobs)
			var id uint64 = 0
			for low := start; low <= end; low += uint16(segmentSize) {
				high := low + uint16(segmentSize) - 1
				if high < low {
					jobs <- Job[uint16]{ID: id, Start: low, End: end}
					break
				}
				jobs <- Job[uint16]{ID: id, Start: low, End: high}
				id++
			}
		}()

		for workerId := 0; workerId < int(workersNum); workerId++ {
			wg.Add(1)

			go func() {
				defer wg.Done()
				worker(jobs, results)
			}()
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		file, err := os.Create(newFilename)
		if err != nil {
			return err
		}
		writer := bufio.NewWriterSize(file, bufferSizeMB<<20)
		defer file.Close()

		var nextID uint64 = 0
		pending := make(map[uint64][]byte)

		for result := range results {
			if result.Err != nil {
				return result.Err
			}

			pending[result.ID] = result.Data

			for {
				data, ok := pending[result.ID]
				if !ok {
					break
				}

				_, err := writer.Write(data)
				if err != nil {
					return err
				}

				delete(pending, nextID)
				nextID++
			}

		}

		err = writer.Flush()
		if err != nil {
			return err
		}
	} else if bits <= 32 {
		start, end := prepareRange[uint32](bits)
		jobs := make(chan Job[uint32])

		go func() {
			defer close(jobs)
			var id uint64 = 0
			for low := start; low <= end; low += uint32(segmentSize) {
				high := low + uint32(segmentSize) - 1
				if high < low {
					jobs <- Job[uint32]{ID: id, Start: low, End: end}
					break
				}
				jobs <- Job[uint32]{ID: id, Start: low, End: high}
				id++
			}
		}()

		for workerId := 0; workerId < int(workersNum); workerId++ {
			wg.Add(1)

			go func() {
				defer wg.Done()
				worker(jobs, results)
			}()
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		file, err := os.Create(newFilename)
		if err != nil {
			return err
		}
		writer := bufio.NewWriterSize(file, bufferSizeMB<<20)
		defer file.Close()

		var nextID uint64 = 0
		pending := make(map[uint64][]byte)

		for result := range results {
			if result.Err != nil {
				return result.Err
			}

			pending[result.ID] = result.Data

			for {
				data, ok := pending[nextID]
				if !ok {
					break
				}

				_, err := writer.Write(data)
				if err != nil {
					return err
				}

				delete(pending, nextID)
				nextID++
			}

		}

		err = writer.Flush()
		if err != nil {
			return err
		}
	} else if bits > 64 {
		return fmt.Errorf("%d cannot be handled!", bits)
	}

	start, end := prepareRange[uint64](bits)
	jobs := make(chan Job[uint64])

	go func() {
		defer close(jobs)
		var id uint64 = 0
		for low := start; low <= end; low += uint64(segmentSize) {
			high := low + uint64(segmentSize) - 1
			if high < low {
				jobs <- Job[uint64]{ID: id, Start: low, End: end}
				break
			}
			jobs <- Job[uint64]{ID: id, Start: low, End: high}
			id++
		}
	}()

	for workerId := 0; workerId < int(workersNum); workerId++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			worker(jobs, results)
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	file, err := os.Create(newFilename)
	if err != nil {
		return err
	}
	writer := bufio.NewWriterSize(file, bufferSizeMB<<20)
	defer file.Close()

	var nextID uint64 = 0
	pending := make(map[uint64][]byte)

	for result := range results {
		if result.Err != nil {
			return result.Err
		}

		pending[result.ID] = result.Data

		for {
			data, ok := pending[nextID]
			if !ok {
				break
			}

			_, err := writer.Write(data)
			if err != nil {
				return err
			}

			delete(pending, nextID)
			nextID++
		}

	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func worker[T Number](jobs <-chan Job[T], results chan<- Result) {
	for job := range jobs {
		primes := genericSegmentedSieve(job.Start, job.End)

		var buf []byte
		for _, p := range primes {
			buf = strconv.AppendUint(buf, uint64(p), 2)
			buf = append(buf, ' ')
			buf = strconv.AppendUint(buf, uint64(p), 10)
			buf = append(buf, '\n')
		}

		results <- Result{
			ID:   job.ID,
			Data: buf,
			Err:  nil,
		}
	}
}

func main() {
	var bits uint8 = 32
	var workersNum uint8 = 0
	var segmentSize int = 0
	var filename string = ""
	// var cacheLine int = cpuid.CPU.CacheLine

	if workersNum == 0 {
		workersNum = uint8(
			min(cpuid.CPU.PhysicalCores, runtime.GOMAXPROCS(0)),
		)
	}

	if segmentSize == 0 {
		segmentSize = cpuid.CPU.Cache.L1D
	}

	run(bits, workersNum, segmentSize, filename, 1)

	// jobs := make(chan Job[uint8])

	// go func() {
	// 	defer close(jobs)
	// 	var id uint64 = 0
	// 	for low := start; low <= end; low += uint8(segmentSize) {
	// 		high := low + uint8(segmentSize) - 1
	// 		if high < low {
	// 			jobs <- Job[uint8]{ID: id, Start: low, End: end}
	// 			break
	// 		}
	// 		jobs <- Job[uint8]{ID: id, Start: low, End: high}
	// 		id++
	// 	}
	// }()

	// for job := range jobs {
	// 	fmt.Println(job)
	// }
}
