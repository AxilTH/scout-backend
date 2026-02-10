// main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("✅ Education Service is alive!"))
	})

	log.Println("🚀 Starting Education Service on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}