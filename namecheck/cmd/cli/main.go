package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jub0bs/namecheck"
	"github.com/jub0bs/namecheck/github"
)

type Result struct {
	Platform  string
	Valid     bool
	Available bool
	Err       error
}

func main() {
	if len(os.Args[1:]) == 0 {
		log.Fatal("username args is required")
	}
	username := os.Args[1]

	var checkers []namecheck.Checker
	const n = 16
	g := &github.GitHub{
		Client: http.DefaultClient,
	}
	for range n {
		checkers = append(checkers, g)
	}
	resultCh := make(chan Result)
	var wg sync.WaitGroup
	for _, checker := range checkers {
		wg.Add(1)
		go check(checker, username, &wg, resultCh)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	var results []Result
	for res := range resultCh {
		results = append(results, res)
	}
	fmt.Println(results)
}

func check(
	checker namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	resultCh chan<- Result,
) {
	defer wg.Done()
	res := Result{
		Platform: checker.String(),
		Valid:    checker.IsValid(username),
	}
	if !res.Valid {
		resultCh <- res
		return
	}
	res.Available, res.Err = checker.IsAvailable(username)
	resultCh <- res
}
