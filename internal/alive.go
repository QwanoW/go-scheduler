package internal

import (
	"fmt"
	"log"
	"net/http"
)

func KeepAlive() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "ALIVE")
		if err != nil {
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
