package lobby

import (
	"encoding/json"
	"fmt"

	"github.com/cjodo/go-server/pkg/socket"
	"github.com/gorilla/websocket"
)

type GameMessage struct {
	Type 	string `json:"type"`
	Code 	string `json:"game-code"`
}

type ChatMessage struct {
	Type			string `json:"chat-message"`
	Game			string `json:"game-code"`
	Message		string `json:"message"`
}

type ReconnectMessage struct {
	Type 			string `json:"type"`
	Game			string `json:"game-code"`
	Board			[9]string `json:"board"`
	Player 		string	`json:"player"`
	IsStarted bool `json:"started"`
	Winner 		string `json:"winner"`
}

func (l *Lobby) handleMessage(message []byte, client *socket.Connection) error {
	type msgType struct {
		Type string `json:"type"`
	}

	var err error

	var msg msgType

	if err := json.Unmarshal(message, &msg); err != nil {
		if websocket.IsCloseError(err) {
			fmt.Println("websocket close: ", err)
		} else {
			return fmt.Errorf("error unmarshalling data: %v\n", err)
		}
	}

	if msg.Type == "" {
		return fmt.Errorf("message had no `type` field")
	}

	switch msg.Type {
	case "chat-message":
		var chatMessage ChatMessage
		if err := json.Unmarshal(message, &chatMessage); err != nil {
			return fmt.Errorf("error unmarshalling chat message: %v\n", err)
		}
		fmt.Println("handle message: ", chatMessage, " from: ", client.Id)

	case "create-game":
		err = l.HandleCreateGame(message, client)
	case "join-game": 
		err = l.HandleJoinGame(message, client)
	case "move-message":
		err = l.HandleMoveMessage(message)
	case "rematch":
		err = l.HandleRematch(message)

	default: 
		return fmt.Errorf("message type not recognized")
	}

	return err
}
