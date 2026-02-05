package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	cc   = flag.Bool("concurrent", false, "If true, run the concurrent function")
)

func sequential(items []string) {
	for _, item := range items {
		fmt.Println("Printing sequential task: " + item)
	}
}

func concurrent(items []string) {
	for _, item := range items {
		go func(item string) {
			fmt.Println("Printing concurrent task: " + item)
		}(item)
	}
}

func main() {
	flag.Parse()
	items := []string{"task1", "task2", "task3", "task4", "task5"}
	
	// Log how long did it take.
	defer func(begin time.Time) {
		fmt.Println("duration:", time.Since(begin))
	}(time.Now())

	if !*cc {
		sequential(items)
	} else {
		concurrent(items)
	}
}