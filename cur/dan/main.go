package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// parseRow is the worker function that processes one CSV row
func parseRow(row string, ch chan<- float64, wg *sync.WaitGroup) {
	defer wg.Done()

	// Simulate real work with delay
	time.Sleep(500 * time.Millisecond)

	// Parse the CSV row manually (split by comma)
	var cost float64
	lastCommaIdx := -1
	for i := len(row) - 1; i >= 0; i-- {
		if row[i] == ',' {
			lastCommaIdx = i
			break
		}
	}

	if lastCommaIdx != -1 {
		costStr := row[lastCommaIdx+1:]
		parsedCost, err := strconv.ParseFloat(costStr, 64)
		if err == nil {
			cost = parsedCost
		}
	}

	// Send the cost to the channel
	ch <- cost
}

func main() {
	startTime := time.Now()

	// Read the CSV file
	file, err := os.Open("../testcur.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Skip header row
	if len(records) > 0 {
		records = records[1:]
	}

	// Create channel and WaitGroup
	costChannel := make(chan float64, len(records))
	var wg sync.WaitGroup

	// Launch a goroutine for each row
	for _, record := range records {
		// Convert record back to string format for parseRow
		rowStr := record[0] + "," + record[1] + "," + record[2]
		wg.Add(1)
		go parseRow(rowStr, costChannel, &wg)
	}

	// Wait for all workers to finish, then close the channel
	go func() {
		wg.Wait()
		close(costChannel)
	}()

	// Collect all costs from the channel
	var totalCost float64
	var rowsProcessed int
	for cost := range costChannel {
		totalCost += cost
		rowsProcessed++
	}

	// Calculate elapsed time
	elapsed := time.Since(startTime)

	// Print results
	fmt.Printf("Total Cost: $%.2f\n", totalCost)
	fmt.Printf("Rows Processed: %d\n", rowsProcessed)
	fmt.Printf("Total Time Taken: %v\n", elapsed)
}
