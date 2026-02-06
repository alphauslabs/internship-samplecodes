package main

// package cur
// 1. The Worker (parseRow):
// Create a function that takes one string (e.g., '1001 and a channel.
// It must split the string and convert the price to a float64.
// The Delay: Add time-Sleep(500 • time-Millisecond) to simulate real work,
// 2. The Foreman (main):
// Launch a Goroutine for every single row in the data slice.
// use a Channel to collect the costs,
// use a sync.WaitGroup to track when all workers are finished.
// 3. The Goal:
// Calculate the TotalCost and RowsProcessed.
// Print the total time taken.

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)


func parseRow(row string, costs chan<- float64, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(500 * time.Millisecond)

	parts := strings.FieldsFunc(row, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t'
	})
	if len(parts) < 2 {
		return
	}

	price, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return
	}

	costs <- price
}

func main() {
	start := time.Now()

	// Read CSV file 
	file, err := os.Open("../testcur.csv")
	if err != nil {
		fmt.Println("Error opening testcur.csv:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Make a data slice of strings 
	var data []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading CSV:", err)
			return
		}

		if len(record) < 3 {
			continue
		}
		//slicing 
		id := strings.TrimSpace(record[0])
		costStr := strings.TrimSpace(record[2])

		if strings.EqualFold(id, "id") || strings.EqualFold(costStr, "cost") {
			continue
		}

		// Send "id,cost"
		data = append(data, id+","+costStr)
	}

	if len(data) == 0 {
		fmt.Println("No valid rows found in testcur.csv.")
		return
	}

	costs := make(chan float64)

	var wg sync.WaitGroup
	wg.Add(len(data))

	for _, row := range data {
		go parseRow(row, costs, &wg)
	}

	go func() {
		wg.Wait()
		close(costs)
	}()

	var totalCost float64
	var rowsProcessed int

	for cost := range costs {
		totalCost += cost
		rowsProcessed++
	}

	elapsed := time.Since(start)

	fmt.Printf("RowsProcessed: %d\n", rowsProcessed)
	fmt.Printf("TotalCost: %.2f\n", totalCost)
	fmt.Printf("TotalTimeTaken: %s\n", elapsed)
}

