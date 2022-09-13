package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Result struct {
	data string
	err  error
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	dc := "Strasbourg"
	results := make(chan Result)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go fetchFrom(ctx, results, dc)
	fmt.Println(<-results)
}

func fetchFrom(ctx context.Context, results chan<- Result, dc string) {
	timer := time.NewTimer(time.Duration(rand.Intn(1000)) * time.Millisecond)
	var res Result
	select {
	case <-ctx.Done():
		fmt.Println("cancelling request to", dc)
		timer.Stop()
		res.err = ctx.Err()
	case <-timer.C:
		res.data = fmt.Sprintf("data from %s", dc)
	}
	results <- res

}
