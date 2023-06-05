package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	// start listening for incoming messages
	go handleMessages()

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// upgrade the GET request to websocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Fatal("Error upgrading the connection to websocket: ", err)
	}

	// register new client
	clients[ws] = true
	log.Println("Client connected")

	// close the connection when this func returns
	defer func() {
		delete(clients, ws)
		ws.Close()
	}()

	for {
		var msg Message
		// read message from websocket
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading the message: %v", err)
			break
		}
		// send the received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// get the next message from the channel
		msg := <-broadcast
		// send message to all connected clients
		for k := range clients {
			err := k.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				k.Close()
				delete(clients, k)
			}
		}
	}
}
