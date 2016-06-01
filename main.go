package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	backend string
	bind    string
)

func init() {
	flag.StringVar(&backend, "backend", "http://127.0.0.1:80", "The backend to proxy to")
	flag.StringVar(&bind, "bind", ":8888", "Address to bind")
}

func main() {
	flag.Parse()
	initAuth()

	u, err := url.Parse(backend)
	if err != nil {
		fmt.Printf("Error parsing backend URL (%v) - %v\n", backend, err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if checkAuth(w, r) {
			proxy.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
	})

	http.ListenAndServe(bind, nil)
}
