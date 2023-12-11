package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gonum/stat/distuv"
)

type Result string

func main() {
	dc := "Strasbourg"
	results := make(chan Result)
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
		fmt.Fprintln(os.Stderr, ctx.Err())
		return
	case res := <-results:
		fmt.Println(res)
	}
}

// Do not modify this function (imagine it's third-party).
func query(ctx context.Context, dc string) (Result, error) {
	latencyMs := 1000 * dist.Rand()
	timer := time.NewTimer(time.Duration(latencyMs) * time.Millisecond)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		fmt.Println("cancelling request to", dc)
		return "", ctx.Err()
	case <-timer.C:
		break
	}
	res := Result(fmt.Sprintf("data from %s", dc))
	return res, nil
}

// see https://en.wikipedia.org/wiki/Log-normal_distribution
var dist = distuv.LogNormal{
	Mu:     0,
	Sigma:  0.5,
	Source: rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
}
