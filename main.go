package main

import (
	"net/http"
)

func setAPI() {
	manager := NewManager()
	http.Handle("/", http.FileServer(http.Dir("./Frontend")))
	http.HandleFunc("/ws", manager.Servews)
}
func main() {
	setAPI()
	http.ListenAndServe(":8080", nil)
}
