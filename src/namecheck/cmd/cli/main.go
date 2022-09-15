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
	var wg sync.WaitGroup
	wg.Add(len(checkers))
	for _, checker := range checkers {
		go check(checker, username, &wg)
	}
	wg.Wait()
}

func check(checker namecheck.Checker, username string, wg *sync.WaitGroup) {
	defer wg.Done()
	if !checker.IsValid(username) {
		fmt.Printf("%q is not valid on %s\n", username, checker.String())
		return
	}
	avail, err := checker.IsAvailable(username)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	if !avail {
		fmt.Printf("%q is valid but unavailable on %s\n", username, checker.String())
		return
	}
	fmt.Printf("%q is valid and available on %s\n", username, checker.String())
}
