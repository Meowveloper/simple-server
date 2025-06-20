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
			Is_Student: true,
		}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(person); err != nil {
			log.Printf("Error encoding json %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
			return
		}
		var new_person Person

		if err := json.NewDecoder(r.Body).Decode(&new_person); err != nil {
			log.Printf("Error decoding JSON Request body : %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if new_person.Name == "" || new_person.Age <= 0 {
			http.Error(w, "Name and age are required fields", http.StatusBadRequest)
			return
		}

		log.Printf("Received new person: Name = %s, Age = %d, City = %s, Is Student? = %t",
			new_person.Name, new_person.Age, new_person.City, new_person.Is_Student,
		)

		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, "Person %s was successfully created.\n", new_person.Name)
	})

	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
