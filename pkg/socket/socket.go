package socket

import (
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
    defer func() {
        fmt.Println("Closing connection for client:", c)
        c.Conn.Close()
    }()

    for {
        message, ok := <-c.Send
        if !ok {
            // Channel closed, exit the loop
            fmt.Println("Send channel closed for client:", c)
            break
        }
        fmt.Println("Sending message:", message)

        if err := c.Conn.WriteJSON(message); err != nil {
            fmt.Printf("Write error for client %s: %v\n", c.Id, err)
            break
        }
    }
}

