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

func parseRow(line string, ch chan float64) {
	parts := strings.Split(line, ",")
	if len(parts) >= 3 {
		price, err := strconv.ParseFloat(parts[2], 64)
		if err == nil {
			time.Sleep(500 * time.Millisecond)
			ch <- price
		}
	}
}

func main() {
	start := time.Now()
	defer func() {
		fmt.Println("Total time taken:", time.Since(start))
	}()

	file, err := os.Open("../testcur.csv")
	if err != nil {
		return
	}
	defer file.Close()

	var data []string
	scanner := bufio.NewScanner(file)
	first := true
	for scanner.Scan() {
		if first {
			first = false
			continue
		}
		data = append(data, scanner.Text())
	}

	ch := make(chan float64)
	var wg sync.WaitGroup

	for _, line := range data {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			parseRow(line, ch)
		}(line)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var TotalCost float64
	RowsProcessed := 0
	for cost := range ch {
		TotalCost += cost
		RowsProcessed++
	}

	fmt.Printf("Total Cost: %.5f\n", TotalCost)
	fmt.Printf("Rows Processed: %d\n", RowsProcessed)
}
