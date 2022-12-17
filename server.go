package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Transaction struct {
	Amount    float64   `json:"amount"`
	TimeStamp time.Time `json:"timestamp"`
}
type Stats struct {
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Max   float64 `json:"max"`
	Min   float64 `json:"min"`
	Count int     `json:"count"`
}
type Location struct {
	City string `json:"city"`
}

var userLocation Location

func setLocation(w http.ResponseWriter, r *http.Request) {
	var l Location
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}
	userLocation = l
	w.WriteHeader(http.StatusNoContent)
}
func resetLocation(w http.ResponseWriter, r *http.Request) {
	userLocation = Location{}
	w.WriteHeader(http.StatusNoContent)
}

var (
	transactions []Transaction
	currentStats Stats
	mut          sync.RWMutex
)

func getstats(w http.ResponseWriter, r *http.Request) {
	if userLocation.City == "" {
		http.Error(w, "Unauthorised User", http.StatusUnauthorized)
		return
	}
	locationHeader := r.Header.Get("Location")
	if locationHeader != userLocation.City {
		http.Error(w, "Unauthorised User", http.StatusUnauthorized)
		return
	}
	mut.Lock()
	stats := currentStats
	mut.Unlock()
	err := json.NewEncoder(w).Encode(stats)
	if err != nil {
		http.Error(w, "Error Encoding Json file", http.StatusInternalServerError)
		return
	}
}
func addTransaction(w http.ResponseWriter, r *http.Request) {
	var t Transaction
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}
	if t.TimeStamp.After(time.Now()) {
		http.Error(w, "Transation is based on future", http.StatusUnprocessableEntity)
	}
	if time.Since(t.TimeStamp) > time.Second*60 {
		http.Error(w, "Transaction is Older than 60 seconds", http.StatusNoContent)
	}
	transactions = append(transactions, t)
	mut.Lock()
	currentStats.Sum += t.Amount
	currentStats.Count++
	if t.Amount > currentStats.Max {
		currentStats.Max = t.Amount
	}
	if t.Amount < currentStats.Min {
		currentStats.Min = t.Amount
	}
	currentStats.Avg = currentStats.Sum / float64(currentStats.Count)
	mut.Unlock()
	w.WriteHeader(http.StatusCreated)
}
func deleteTransactions(w http.ResponseWriter, r *http.Request) {
	transactions = nil
	mut.Lock()
	currentStats = Stats{}
	mut.Unlock()
	w.WriteHeader(http.StatusNoContent)
}
func main() {
	r := mux.NewRouter()
	// Handler Functions
	r.HandleFunc("/stats", getstats).Methods("GET")
	r.HandleFunc("/transaction", addTransaction).Methods("POST")
	r.HandleFunc("/transaction", deleteTransactions).Methods("DELETE")
	r.HandleFunc("/location", setLocation).Methods("POST")
	r.HandleFunc("/location", resetLocation).Methods("DELETE")
	//Server is running on 8080 Port
	fmt.Println("Server is running at Port at 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
