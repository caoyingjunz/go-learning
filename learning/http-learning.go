package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("http " + httpVersion + " ok\n"))
	})

	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}

var (
	httpVersion string
	httpPort    string
)

func init() {
	flag.StringVar(&httpVersion, "httpVersion", "", "httpVersion for test.")
	flag.StringVar(&httpPort, "httpPort", "80", "httpPort for test.")
}
