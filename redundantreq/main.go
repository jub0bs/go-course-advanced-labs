package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gonum/stat/distuv"
)

type Result string

func main() {
	dc := "Strasbourg"
	results := make(chan Result, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()
	go func() {
		res, err := query(ctx, dc)
		if err != nil {
			return
		}
		results <- res
	}()
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Fprintln(os.Stderr, ctx.Err())
			return
		}
	case res := <-results:
		fmt.Println(res)
	}
}

// Do not modify the file beyond this point.

func query(ctx context.Context, dc string) (Result, error) {
	mu.Lock()
	latencyMs := 1000 * dist.Rand()
	mu.Unlock()
	timer := time.NewTimer(time.Duration(latencyMs) * time.Millisecond)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-timer.C:
		res := Result(fmt.Sprintf("received response from %s", dc))
		return res, nil
	}
}

var (
	// see https://en.wikipedia.org/wiki/Log-normal_distribution
	dist = distuv.LogNormal{
		Mu:     0,
		Sigma:  0.5,
		Source: rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
	}
	mu sync.Mutex
)
