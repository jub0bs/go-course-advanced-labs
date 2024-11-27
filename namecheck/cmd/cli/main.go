package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jub0bs/namecheck/github"
)

func main() {
	if len(os.Args[1:]) == 0 {
		log.Fatal("username args is required")
	}
	username := os.Args[1]

	var checkers []Checker
	const n = 4
	g := &github.GitHub{
		Client: http.DefaultClient,
	}
	for range n {
		checkers = append(checkers, g)
	}
	for _, checker := range checkers {
		check(checker, username)
	}
}

func check(checker Checker, username string) {
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
