package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	fmt.Println("hello")
}
