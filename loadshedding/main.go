package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	sem chan struct{}
}

func NewServer(maxConcurrent int) (*Server, error) {
	if maxConcurrent <= 0 {
		return nil, errors.New("nonpositive maxConcurrent")
	}
	srv := Server{
		sem: make(chan struct{}, maxConcurrent),
	}
	return &srv, nil
}

func (srv *Server) HelloWorld() http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		srv.sem <- struct{}{}
		defer func() { <-srv.sem }()
		time.Sleep(5 * time.Second) // simulate work
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "Hello, World!")
	}
	return h
}

func main() {
	const maxConcurrent = 2
	srv, err := NewServer(maxConcurrent)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", srv.HelloWorld())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
