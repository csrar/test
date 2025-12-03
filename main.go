package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Event struct {
	ID          uuid.UUID `json:""`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
}

var db *sql.DB

func createEvent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID          uuid.UUID `json:"id"`
		Title       string    `json:"title"`
		Description *string   `json:"description"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		CreatedAt   time.Time `json:"created_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Title == "" || len(input.Title) > 100 || !input.StartTime.Before(input.EndTime) {
		http.Error(w, "validation failed", http.StatusBadRequest)
		return
	}

	id := uuid.New()
	createdAt := time.Now().UTC()

	_, err := db.ExecContext(r.Context(),
		`INSERT INTO events (id, title, description, start_time, end_time, created_at) values($1, $2, $3, $4, $5, $6)`,
		id, input.Title, input.Description, input.StartTime, input.EndTime, createdAt)
	if err != nil {
		http.Error(w, "Error inserting", http.StatusInternalServerError)
	}

	event := Event{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		CreatedAt:   createdAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func main() {

	var err error
	db, err := sql.Open("postgres", "postgresql://postgres:xxx@db.jrebekphkfviwzxsnctu.supabase.co:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/events", createEvent).Methods("POST")
	// r.HandleFunc("/events", listEvents).Methods("GET")
	// r.HandleFunc("/events/{id}", getEventById).Methods("GET")
	fmt.Println("app started")
	http.ListenAndServe(":8080", r)

}
