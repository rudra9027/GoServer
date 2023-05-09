package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	WebsocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	client ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		client: make(ClientList),
	}
}
func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.client[client] = true
}
func (m *Manager) remove(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.client[client]; ok {
		client.connection.Close()
		delete(m.client, client)
	}
}
func (m *Manager) Servews(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")
	conn, err := WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m)
	m.addClient(client)
	go client.ReadMessages()
	go client.WriteMessage()

}
func (m *Manager) RemoveClient(c *Client) {

}
