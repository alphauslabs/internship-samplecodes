package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func processLine(line []string, result chan<- float64) {
	time.Sleep(500 * time.Millisecond)

	price := line[2]
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println(err)
		result <- 0.0
		return
	}
	result <- priceFloat
}

func main() {
	startTime := time.Now()
	fileCSV, err := os.Open("../testcur.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	read := csv.NewReader(fileCSV)
	lines, err := read.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	total := make(chan float64)
	var wg sync.WaitGroup

	rows := lines[1:]

	for _, row := range rows {
		wg.Add(1)

		go func(r []string) {
			defer wg.Done()
			processLine(r, total)
		}(row)
	}

	go func() {
		wg.Wait()
		close(total)
	}()

	var totals float64
	for price := range total {
		totals += price
	}

	fmt.Printf("Total: %f\n", totals)
	fmt.Printf("Rows processed: %d\n", len(rows))
	fmt.Printf("Duration: %s\n", time.Since(startTime))
}
