package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var origins = "*"

func main() {
	http.HandleFunc("/events", eventsHandler)
	fmt.Println("SSE server is running on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Unable to start SSE server: %s", err.Error())
	}
}

// SSE handler
func eventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", origins)
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	// Channel for client disconnection
	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)

	timer := time.NewTicker(time.Second * 5)
	defer timer.Stop()

	for {
		select {
		case <-clientGone:
			fmt.Println("Client disconnected")
			return
		case <-timer.C:
			// Create message body
			fmt.Fprintf(w, "event:ticker\ndata: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
			// Flush the response
			err := rc.Flush()
			if err != nil {
				fmt.Println("Unable to flush response")
				return
			}
		}
	}

}
