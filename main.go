package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "addr to bind on")
}

func main() {
	wh, rh := startServer()
	http.Handle("/write", wh)
	http.Handle("/read", rh)
	log.Fatal(http.ListenAndServe(addr, nil))
}
