package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Person struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	City       string `json:"city,omitempty"`
	Is_Student bool   `json:"is_student"`
}

type Error_Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", handle_root)
	mux.HandleFunc("/hello", handle_hello)
	mux.HandleFunc("/greet/", handle_greet)
	mux.HandleFunc("/api/person", handle_get_person)
	mux.HandleFunc("/api/register", handle_register)
	mux.HandleFunc("/api/search", handle_search)
	mux.HandleFunc("/api/form-submit", handle_form_submit)

	fs := http.FileServer(http.Dir("static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", logging_middleware(mux)))
}

func handle_root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "static/index.html")
		return
	}
	send_JSON_error(w, "not found", http.StatusNotFound, "the requested resource was not found")
}

func handle_hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello there!")
}

func handle_greet(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/greet/")

	if name == "" {
		name = "Guest"
	}

	fmt.Fprintf(w, "Greetings %s!", name)
}

func handle_get_person(w http.ResponseWriter, r *http.Request) {
	person := Person{
		Name:       "Alice",
		Age:        24,
		Is_Student: true,
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(person); err != nil {
		log.Printf("Error encoding json %v", err)
		send_JSON_error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func handle_register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		return
	}
	var new_person Person

	if err := json.NewDecoder(r.Body).Decode(&new_person); err != nil {
		log.Printf("Error decoding JSON Request body : %v", err)
		send_JSON_error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if new_person.Name == "" || new_person.Age <= 0 {
		send_JSON_error(w, "Name and age are required fields", http.StatusBadRequest)
		return
	}

	log.Printf("Received new person: Name = %s, Age = %d, City = %s, Is Student? = %t",
		new_person.Name, new_person.Age, new_person.City, new_person.Is_Student,
	)

	w.WriteHeader(http.StatusCreated)

	fmt.Fprintf(w, "Person %s was successfully created.\n", new_person.Name)
}

func handle_search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	limit := r.URL.Query().Get("limit")

	if query == "" {
		send_JSON_error(w, "Missing query string", http.StatusBadRequest)
		return
	}

	limit_value := 0
	if limit != "" {
		var err error
		limit_value, err = strconv.Atoi(limit)
		if err != nil {
			send_JSON_error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "searching for: '%s' with a limit of %d results", query, limit_value)
}

func handle_form_submit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		send_JSON_error(w, "request method not allowed", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form data %v", err)
		send_JSON_error(w, "error parsing form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if username == "" || password == "" {
		send_JSON_error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	log.Printf("received form submission: username=%s, email=%s.", username, email)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "form submitted successfully! welcome, %s", username)
}

func send_JSON_error(w http.ResponseWriter, message string, status_code int, details ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)

	err_response := Error_Response{
		Message: message,
		Code:    status_code,
	}

	if len(details) > 0 && details[0] != "" {
		err_response.Details = details[0]
	}

	if err := json.NewEncoder(w).Encode(err_response); err != nil {
		log.Printf("failed to write error response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

type logging_response_writer struct {
	http.ResponseWriter
	status_code int
}

func new_logging_response_writer(w http.ResponseWriter) *logging_response_writer {
	return &logging_response_writer{w, http.StatusOK}
}
func (lrw *logging_response_writer) Write_Header(code int) {
	lrw.status_code = code
	lrw.ResponseWriter.WriteHeader(code)
}

func logging_middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := new_logging_response_writer(w) // pointer

		next.ServeHTTP(lrw, r)

		log.Printf("%s %s %s %d %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			lrw.status_code,
			time.Since(start),
		)
	})
}
