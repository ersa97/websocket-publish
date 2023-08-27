package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type websocketConnection struct {
	Conn *websocket.Conn
	Key  string
}

var WsConnectionPool []*websocketConnection

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn.SetReadDeadline(time.Time{})
	conn.SetWriteDeadline(time.Time{})
	key := r.URL.Query().Get("key")
	ws := websocketConnection{
		Conn: conn,
		Key:  key,
	}
	WsConnectionPool = append(WsConnectionPool, &ws)
	log.Println("new websocket channel has been created with key :", key)

}
func publish(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	var conn *websocket.Conn
	for _, v := range WsConnectionPool {
		if v.Key == key {
			conn = v.Conn
		}
	}
	if conn == nil {
		log.Println("connection not found")
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte("halo from server key "+key))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("new message has been sent to key :", key)
	r.Response = &http.Response{
		StatusCode: http.StatusOK,
	}
}
func main() {
	http.HandleFunc("/ws", ws)
	http.HandleFunc("/publish", publish)

	log.Println("Server starting at :1322")
	http.ListenAndServe(":1322", nil)
}
