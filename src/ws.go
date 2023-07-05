package src

import (
	"log"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (pool *Pool) Start() {
	manager := GetManager()
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Println("Size of connection Pool:", len(pool.Clients))
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.Println("Size of connection Pool:", len(pool.Clients))
			if len(pool.Clients) == 0 {
				log.Println("reach empty poll")
				delete(manager.Pools, client.ExternalUUID)
			}
			break
		case msg := <-pool.Broadcast:
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(msg); err != nil {
					log.Println("caught error: ", err)
					pool.Unregister <- client
					return
				}
			}
		}
	}
}
