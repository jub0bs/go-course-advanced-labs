package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jub0bs/namecheck"
	"github.com/jub0bs/namecheck/github"
	"github.com/jub0bs/namecheck/twitter"
)

type Status int

type Result struct {
	Username  string
	Platform  string
	Valid     bool
	Available bool
}

const (
	Unknown Status = iota
	Active
	Suspended
	Available
)

func main() {
	if len(os.Args[1:]) == 0 {
		log.Fatal("username args is required")
	}
	username := os.Args[1]

	var checkers []namecheck.Checker
	for i := 0; i < 3; i++ {
		t := &twitter.Twitter{
			Client: http.DefaultClient,
		}
		g := &github.GitHub{
			Client: http.DefaultClient,
		}
		checkers = append(checkers, t, g)
	}
	results := make(chan Result, len(checkers))
	errc := make(chan error, len(checkers))
	var wg sync.WaitGroup
	wg.Add(len(checkers))
	for _, checker := range checkers {
		go check(checker, username, &wg, results, errc)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for {
		select {
		case err := <-errc:
			const tmpl = "namecheck: some error occurred: %s\n"
			fmt.Fprintf(os.Stderr, tmpl, err)
			return
		case res, ok := <-results:
			if !ok {
				return
			}
			fmt.Println(res)
		}
	}
}

func check(
	checker namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	results chan<- Result,
	errc chan<- error,
) {
	defer wg.Done()
	res := Result{
		Username: username,
		Platform: checker.String(),
	}
	res.Valid = checker.IsValid(username)
	if !res.Valid {
		results <- res
		return
	}
	avail, err := checker.IsAvailable(username)
	if err != nil {
		errc <- err
		return
	}
	res.Available = avail
	results <- res
}
