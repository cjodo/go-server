package socket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Id			string
	Token 	string
	Player 	string
	Conn		*websocket.Conn
	Send		chan interface{}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Connection) HandleWrite() {
	for {
		message := <-c.Send

		json, err := json.Marshal(message)
		if err != nil {
			fmt.Println("error marshalling json: ", err)
			break
		}

		if err := c.Conn.WriteMessage(1, json); err != nil {
			fmt.Println(err)
			break
		}
	}
}
