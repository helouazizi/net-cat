package main

import (
	"fmt"
	"sync"
	"time"
)

func printMessage(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Decrease the counter when the goroutine finishes

	for i := 0; i < 3; i++ {
		fmt.Printf("Goroutine %d - Message %d\n", id, i+1)
		time.Sleep(500 * time.Millisecond) // Sleep for half a second
	}
}

func main() {
	var wg sync.WaitGroup

	// Start multiple goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Increment the counter for each goroutine
		go printMessage(i, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("All goroutines completed")
}
