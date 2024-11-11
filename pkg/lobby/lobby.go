package lobby

import (
	"fmt"
	"sync"

	"github.com/cjodo/go-server/pkg/game"
	"github.com/cjodo/go-server/pkg/socket"
	"github.com/gorilla/websocket"
)

type Lobby struct {
	clients				map[string]*socket.Connection // using the connection.Id as token
	Broadcast			chan interface{}
	Register			chan *socket.Connection
	Unregister		chan *socket.Connection
	Games					map[string]*game.Game 
	mu 						sync.Mutex
}

func NewLobby() *Lobby {
	return &Lobby{
		clients:    make(map[string]*socket.Connection),
		Broadcast:  make(chan interface{}),
		Register:   make(chan *socket.Connection),
		Unregister: make(chan *socket.Connection),
		Games:			make(map[string]*game.Game),		
		mu: 				sync.Mutex{},			
	}
}

func (l *Lobby) Run() {
	for {
		select {
		case msg := <-l.Broadcast:
			// Unused
			fmt.Println("Broadcast: ", msg)
		case client := <-l.Register:
			// ping and clean up all nil connections

			if _, exists := l.clients[client.Id]; exists {
				go l.handleReconnect(client)
				go client.HandleWrite()
				go l.handleClientRead(client)
				break
			} else {
				go l.registerConn(client)
				go client.HandleWrite()
				go l.handleClientRead(client)
			}

			l.mu.Lock()
			l.pingClients()
			l.mu.Unlock()

		case client := <-l.Unregister:
			client.Conn.Close()
		}
	}
}

func (l *Lobby) registerConn(client *socket.Connection) {
	l.clients[client.Id] = client
}

func (l * Lobby) handleClientRead(client *socket.Connection) {
	go func () {
		defer client.Conn.Close()

		for {
			_, p, err := client.Conn.ReadMessage()
			if err != nil {
				fmt.Println("error reading message: ", err)
				l.Unregister <- client
				break
			}

			if err = l.handleMessage(p, client); err != nil {
				fmt.Println(err)
			}

		}
	}()
}

func (l *Lobby) pingClients () {
	fmt.Println("Pinging Clients")

	for _, client := range l.clients {
		// Send a ping to the client
		err := client.Conn.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			fmt.Printf("Client disconnected: %v\n", err)
			delete(l.clients, client.Id)
			continue
		}
	}
}
