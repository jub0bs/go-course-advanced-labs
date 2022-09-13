package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gorilla/mux"
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
	Error     error
}

const (
	Unknown Status = iota
	Active
	Suspended
	Available
)

var (
	visits uint64
	m      = make(map[string]uint)
	mu     sync.Mutex
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/check", handleCheck)
	r.HandleFunc("/visits", handleVisits)
	r.HandleFunc("/details", handleDetails)
	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dec := json.NewEncoder(w)
	mu.Lock()
	defer mu.Unlock()
	if err := dec.Encode(m); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func handleVisits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	entity := struct {
		Count uint64 `json:"visits"`
	}{
		Count: atomic.LoadUint64(&visits),
	}
	dec := json.NewEncoder(w)
	if err := dec.Encode(entity); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func handleCheck(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&visits, 1)
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "'username' query param is required", http.StatusBadRequest)
		return
	}
	mu.Lock()
	m[username]++
	mu.Unlock()
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
	results := make(chan Result)
	var wg sync.WaitGroup
	wg.Add(len(checkers))
	for _, checker := range checkers {
		go check(checker, username, &wg, results)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	w.Header().Set("Content-Type", "application/json")
	type jsonResult struct {
		Platform  string `json:"platform"`
		Valid     string `json:"valid"`
		Available string `json:"available"`
	}
	jsonResults := make([]jsonResult, 0, len(checkers))
	for result := range results {
		res := jsonResult{
			Platform:  result.Platform,
			Valid:     fmt.Sprintf("%t", result.Valid),
			Available: availabilityStatus(result),
		}
		jsonResults = append(jsonResults, res)
	}
	entity := struct {
		Username string       `json:"username"`
		Results  []jsonResult `json:"results"`
	}{
		Username: username,
		Results:  jsonResults,
	}
	dec := json.NewEncoder(w)
	if err := dec.Encode(entity); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func check(
	checker namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	results chan<- Result,
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
	res.Available = avail
	if err != nil {
		res.Error = err
	}
	results <- res
}

func availabilityStatus(res Result) string {
	if res.Error != nil {
		return "unknown"
	}
	return fmt.Sprintf("%t", res.Available)
}
