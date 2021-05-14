package main

import (
	"flag"
	"fmt"
	"github.com/Sahil624/websocket_channels"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var groups = []string{"Mango", "Orange", "Apple"}

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(channel *websocket_channels.Channel,w http.ResponseWriter, r *http.Request) {
	groupName := groups[rand.Intn(len(groups))]
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	channel.GroupAdd(conn, groupName)
	channel.GroupSend(groupName, "Someone Entered This Group")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		channel.GroupSend(groupName, string(message))
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}


func main() {
	channel := websocket_channels.New(websocket_channels.ChannelsConfig{
		Debug: true,
	})
	rand.Seed(time.Now().Unix())
	fmt.Println("Channel", channel)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(channel, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
