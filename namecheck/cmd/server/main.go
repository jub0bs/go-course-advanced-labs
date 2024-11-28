package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"sync"
	"time"

	"github.com/jub0bs/cors"
	"github.com/jub0bs/namecheck"
	"github.com/jub0bs/namecheck/github"
)

type Result struct {
	Platform  string `json:"platform"`
	Valid     bool   `json:"valid"`
	Available bool   `json:"available"`
	Err       error  `json:"error,omitempty"`
}

var (
	stats = make(map[string]uint)
	mu    sync.Mutex
)

func main() {
	// create a custom HTTP request multiplexer and register your handler to pattern GET /hello
	mux := http.NewServeMux()
	mux.HandleFunc("GET /check", handleCheck)
	mux.HandleFunc("GET /stats", handleStats)

	// instantiate a CORS middleware whose config suits your needs
	corsMw, err := cors.NewMiddleware(cors.Config{
		Origins: []string{"https://jub0bs.github.io"},
		ExtraConfig: cors.ExtraConfig{
			PrivateNetworkAccess: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// apply your CORS middleware to your HTTP request multiplexer
	handler := corsMw.Wrap(mux)

	// start the server on port 8080; make sure to use your custom handler
	if err := http.ListenAndServe(":8080", handler); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	mu.Lock()
	statsCopy := maps.Clone(stats)
	mu.Unlock()
	if err := enc.Encode(statsCopy); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleCheck(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mu.Lock()
	stats[username]++
	mu.Unlock()
	var checkers []namecheck.Checker
	const n = 16
	g := &github.GitHub{
		Client: http.DefaultClient,
	}
	for range n {
		checkers = append(checkers, g)
	}
	resultCh := make(chan Result)
	errorCh := make(chan error)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	for _, checker := range checkers {
		wg.Add(1)
		go check(ctx, checker, username, &wg, resultCh, errorCh)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	var results []Result
	var finished bool
	for !finished {
		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusInternalServerError)
			return
		case <-errorCh:
			cancel()
			fmt.Println(ctx.Err())
			w.WriteHeader(http.StatusInternalServerError)
			return
		case res, ok := <-resultCh:
			if !ok {
				finished = true
				continue
			}
			results = append(results, res)
		}
	}
	data := struct {
		Username string   `json:"username"`
		Results  []Result `json:"results,omitempty"`
	}{
		Username: username,
		Results:  results,
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func check(
	ctx context.Context,
	checker namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	resultCh chan<- Result,
	errorCh chan<- error,
) {
	defer wg.Done()
	res := Result{
		Platform: checker.String(),
		Valid:    checker.IsValid(username),
	}
	if !res.Valid {
		send(ctx, resultCh, res)
		return
	}
	avail, err := checker.IsAvailable(ctx, username)
	if err != nil {
		send(ctx, errorCh, err)
		return
	}
	res.Available = avail
	send(ctx, resultCh, res)
}

func send[T any](ctx context.Context, ch chan<- T, v T) {
	select {
	case <-ctx.Done():
		return
	case ch <- v:
		return
	}
}
