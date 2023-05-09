package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	//egress to avoid the concurrent write to websocket connection
	egress chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}
func (c *Client) ReadMessages() {
	defer func() {
		c.manager.RemoveClient(c)
	}()
	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Println(messageType)
		log.Println(string(payload))
	}
}
func (c *Client) WriteMessage() {
	defer func() {
		c.manager.RemoveClient(c)
	}()
	for {
		select {
		case message, ok := <-c.egress:
			if !ok {

			}
		}
	}
}
