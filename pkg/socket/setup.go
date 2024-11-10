package socket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// What to setup :-()
func SetupHandlers(conn *websocket.Conn) {
	conn.SetPingHandler(func(appData string) error {
		fmt.Printf("Ping received %v\n", appData)
		return nil
	})

}
