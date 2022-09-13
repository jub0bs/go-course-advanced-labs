package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct{}

func (s *Server) Foo() http.HandlerFunc {
	msg := s.initializeFoo()
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, msg)
	}
	return f
}

func (*Server) initializeFoo() string {
	time.Sleep(30 * time.Second) // simulate work
	return "Foo!"
}

func (*Server) Bar() http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "Bar!")
	}
	return f
}

func main() {
	var srv Server
	http.Handle("/foo", srv.Foo())
	http.Handle("/bar", srv.Bar())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
