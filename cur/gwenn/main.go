package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Result holds parsed cost info
type Result struct {
	cost float64
}

// Worker function (Requirement #1)
func parseRow(row string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	// Simulate slow database
	time.Sleep(100 * time.Millisecond)

	fields := strings.Split(row, ",")
	if len(fields) < 3 {
		return
	}

	cost, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return
	}

	results <- Result{cost: cost}
}

func main() {
	start := time.Now()

	file, err := os.Open("/home/gwenn/internship-samplecodes/cur/testcur.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Skip header
	scanner.Scan()

	results := make(chan Result)
	var wg sync.WaitGroup

	rowsProcessed := 0
	totalCost := 0.0

	// Collector goroutine
	go func() {
		for res := range results {
			totalCost += res.cost
			rowsProcessed++
		}
	}()

	// Read file line by line
	for scanner.Scan() {
		line := scanner.Text()

		wg.Add(1)
		go parseRow(line, results, &wg)
	}

	// Wait for workers
	wg.Wait()
	close(results)

	elapsed := time.Since(start)

	fmt.Println("Rows processed:", rowsProcessed)
	fmt.Printf("Total cost: %.2f\n", totalCost)
	fmt.Println("Time taken:", elapsed)
}
