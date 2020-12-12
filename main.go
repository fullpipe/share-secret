package main

import (
	"io"
	"net/http"
	"net/rpc"

	"github.com/rs/cors"
)

func main() {
	ns := NewSecret()
	rpc.Register(ns)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		res := NewRPCRequest(r.Body).Call()
		io.Copy(w, res)
	})

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}
