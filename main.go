package main

import (
	"fmt"
	"net/http"
	"os"
)

type server struct {
	mux *http.ServeMux
}

func newServer() (*server, error) {
	srv := &server{
		mux: http.NewServeMux(),
	}

	srv.mux.HandleFunc("/", srv.indexHandler())

	return srv, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		ip := r.Header.Get("X-Forwarded-For")

		if ip == "" {
			ip = r.RemoteAddr
		}

		fmt.Fprintf(w, ip)
	}
}

func run(args []string) error {
	port := 8080
	addr := fmt.Sprintf("0.0.0.0:%d", port)

	srv, err := newServer()
	if err != nil {
		return err
	}

	fmt.Printf("server listening on :%d\n", port)

	return http.ListenAndServe(addr, srv)
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
