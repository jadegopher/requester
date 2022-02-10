package test

import (
	"net/http"
	"time"
)

func FastAndSlowHandlers(delay time.Duration) http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/fast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("a"))
	})

	r.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("a"))
	})

	return r
}
