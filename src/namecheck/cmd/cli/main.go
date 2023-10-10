package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	for _, checker := range checkers {
		check(checker, username)
	}
}

func check(checker namecheck.Checker, username string) {
	if !checker.IsValid(username) {
		fmt.Printf("%q is not valid on %s\n", username, checker.String())
		return
	}
	avail, err := checker.IsAvailable(username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if !avail {
		fmt.Printf("%q is valid but unavailable on %s\n", username, checker.String())
		return
	}
	fmt.Printf("%q is valid and available on %s\n", username, checker.String())
}
