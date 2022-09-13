package main

import (
	"fmt"
	"sync"
)

func main() {
	printTenIntsConcurrently()
}

func printTenIntsConcurrently() {
	var wg sync.WaitGroup
	const n = 10
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(j int) {
			defer wg.Done()
			fmt.Println(j)
		}(i)
	}
	wg.Wait()
}
