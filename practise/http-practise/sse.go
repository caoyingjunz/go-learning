package main

import (
	"fmt"
	"net/http"
	"time"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {

		fmt.Println("start time", time.Now())
		_, err := fmt.Fprintf(w, "data: %s", "hello world")
		if err != nil {
			fmt.Println("err", err)
		}

		flusher, _ := w.(http.Flusher)
		flusher.Flush()

		fmt.Println("end time", time.Now())
		fmt.Println("==========")

		<-time.After(5 * time.Second)
	}
}

func main() {
	http.HandleFunc("/sse", sseHandler)
	http.ListenAndServe(":8081", nil)
}
