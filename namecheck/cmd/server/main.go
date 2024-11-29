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
	"github.com/jub0bs/errutil"
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

var sema = make(chan struct{}, 8)

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
	results := make([]Result, len(checkers))
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	errCh := make(chan error)
	var wg sync.WaitGroup
	for i, checker := range checkers {
		sema <- struct{}{}
		wg.Add(1)
		go func() {
			res, err := check(ctx, checker, username, &wg)
			if err != nil {
				cancel()
				errCh <- err
				return
			}
			results[i] = *res
		}()
	}
	wg.Wait()
	select {
	case err := <-errCh:
		if errua, ok := errutil.Find[*namecheck.UnknownAvailabilityError](err); ok {
			fmt.Println(errua.Platform, errua.Username)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	default:
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
) (*Result, error) {
	defer wg.Done()
	defer func() { <-sema }()
	res := Result{
		Platform: checker.String(),
		Valid:    checker.IsValid(username),
	}
	if !res.Valid {
		return &res, nil
	}
	avail, err := checker.IsAvailable(ctx, username)
	if err != nil {
		return nil, err
	}
	res.Available = avail
	return &res, nil
}
