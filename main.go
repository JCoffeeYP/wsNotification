package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
	"websocketsNotification/src"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var wsConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	externalUUID := strings.Split(url, "/")[2] // todo придумать как избавиться от индексов

	c, err := wsConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Caught error during upgrade ws connection:", err)
		return
	}

	manager := src.GetManager()

	pool, ok := manager.Pools[externalUUID]
	if !ok {
		pool = *src.NewPool()
	}

	client := &src.Client{
		ExternalUUID: externalUUID,
		Pool:         &pool,
		Conn:         c,
	}
	go pool.Start()
	manager.Pools[externalUUID] = pool
	pool.Register <- client
	client.Read()
}

func genChMsg() {
	manager := src.GetManager()
	for {
		for _, value := range manager.Pools {
			msg := src.Message{
				Type: websocket.TextMessage,
				Body: time.Now().String(),
			}
			value.Broadcast <- msg
		}
		time.Sleep(2 * time.Second)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	r := mux.NewRouter()
	r.HandleFunc("/ws/{uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", wsHandler)
	go genChMsg()
	log.Fatal(http.ListenAndServe(*addr, r))
}
