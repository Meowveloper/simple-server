package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Person struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	City       string `json:"city,omitempty"`
	Is_Student bool   `json:"is_student"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to home page")
	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello there!")
	})

	http.HandleFunc("/greet/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/greet/")

		if name == "" {
			name = "Guest"
		}

		fmt.Fprintf(w, "Greetings %s!", name)
	})

	http.HandleFunc("/api/person", func(w http.ResponseWriter, r *http.Request) {
		person := Person{
			Name:       "Alice",
			Age:        24,
			City:       "New York",
			Is_Student: true,
		}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(person); err != nil {
			log.Printf("Error encoding json %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
