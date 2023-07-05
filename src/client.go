package src

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ExternalUUID string
	Conn         *websocket.Conn
	Pool         *Pool
}

func (client *Client) Read() {
	defer func() {
		client.Pool.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
	}
}
