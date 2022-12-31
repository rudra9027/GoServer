package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Location struct {
	City string `json:"city"`
}
type Stats struct {
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Max   float64 `json:max"`
	Min   float64 `json"min"`
	Count int     `json:"count"`
}
type Transaction struct {
	Amount    float64   `json"amount"`
	Timestamp time.Time `json:"timestamp"`
}

const buffersize = 1000

var (
	buffer       [buffersize]Transaction
	head         int
	tail         int
	count        int
	CurrentStats Stats
	location     Location
)

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	var t Transaction
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	if t.Timestamp.After(now) {
		http.Error(w, "Transaction is outdated", http.StatusUnprocessableEntity)
		return
	}
	if now.Sub(t.Timestamp) > 60*time.Second {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	buffer[head] = t
	head = (head + 1) % buffersize
	if count < buffersize {
		count++
	} else {
		tail = (tail + 1) % buffersize
	}
	CurrentStats.Sum += t.Amount
	CurrentStats.Count++

	if CurrentStats.Max > t.Amount {
		CurrentStats.Max = t.Amount
	}
	if CurrentStats.Min < t.Amount {
		CurrentStats.Min = t.Amount
	}
	CurrentStats.Avg = CurrentStats.Sum / float64(CurrentStats.Count)
	w.WriteHeader(http.StatusCreated)
}
func handleStatistics(w http.ResponseWriter, r *http.Request) {

	if location.City != "" && r.Header.Get("location") != location.City {
		http.Error(w, "unauthorized.", http.StatusUnauthorized)
		return
	}

	var stats Stats
	//stats = CurrentStats

	now := time.Now().UTC().Unix()

	for i := tail; i != head; i = (i + 1) % buffersize {
		t := buffer[i]
		if t.Timestamp.Unix() > now-60 {
			stats.Sum += t.Amount
			stats.Count++
			if t.Amount > stats.Max {
				stats.Max = t.Amount
			}
			if stats.Min == 0 || t.Amount < stats.Min || stats.Count == 1 {
				stats.Min = t.Amount
			}
		}
	}
	stats.Avg = stats.Sum / float64(stats.Count)
	w.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(stats)
	if err != nil {
		http.Error(w, "Get Query is not possible", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func set_location(w http.ResponseWriter, r *http.Request) {
	var loc Location
	err := json.NewDecoder(r.Body).Decode(&loc)
	if err != nil {
		http.Error(w, "Invalid Json to parse", http.StatusBadRequest)
		return
	}
	if loc.City == "" {
		http.Error(w, "City is Empty", http.StatusBadRequest)
		return
	}
	location = loc
	w.WriteHeader(http.StatusOK)
}
func handleResetLocation(w http.ResponseWriter, r *http.Request) {
	location = Location{}
	w.WriteHeader(http.StatusOK)
}

func deleteTransaction(w http.ResponseWriter, r *http.Request) {
	for i := tail; i != head; i = (i + 1) % buffersize {
		buffer[i] = Transaction{}
	}
	head = 0
	tail = 0
	count = 0
	CurrentStats = Stats{}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()
	fmt.Println("Hello World")
	r.HandleFunc("/transactions", AddTransaction)
	r.HandleFunc("/statistics", handleStatistics)
	r.HandleFunc("/delete", deleteTransaction)
	r.HandleFunc("/set_location", set_location)
	http.ListenAndServe(":8080", r)
}
