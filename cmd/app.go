package main

import (
	"fmt"
	"log"
	"net/http"
	"private/scheduler/internal"
)

func main() {
	// keep app alive
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "ALIVE")
		if err != nil {
			return
		}
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// start bot
	internal.StartBot()
}
