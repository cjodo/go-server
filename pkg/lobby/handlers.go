package lobby

import (
	"encoding/json"
	"fmt"

	"github.com/cjodo/go-server/pkg/game"
	"github.com/cjodo/go-server/pkg/socket"
)

func (l *Lobby) HandleJoinGame (message []byte, client *socket.Connection) error { 
	var joinGame GameMessage
	if err := json.Unmarshal(message, &joinGame); err != nil {
		return fmt.Errorf("error unmarshalling join-game message: %v\n", err)
	}

	targetGame, exists := l.Games[joinGame.Code]
	if !exists {
		return fmt.Errorf("game doesn't exist at: %v", joinGame.Code)
	}
	// TODO destroy inactive games
	if _, exists := targetGame.Players[client.Id]; exists {
		return fmt.Errorf("client already in game: %v, %v", targetGame.Code, client.Id)
	}
	//client doesnt exist
	targetGame.Players[client.Id] = client

	if len(targetGame.Players) == 2 {
		targetGame.Start(false)
		return nil
	}

	return nil
}

func (l *Lobby) HandleMoveMessage (message []byte) error {
	var moveMessage game.MoveMessage
	if err := json.Unmarshal(message, &moveMessage); err != nil {
		return fmt.Errorf("error unmarshalling move-message: %v\n", err)
	}

	targetGame, exists := l.Games[moveMessage.Code]
	if !exists {
		return fmt.Errorf("game doesn't exist at: %v", moveMessage.Code)
	}

	if targetGame.Winner != "" {
		return fmt.Errorf("winner on board %v", targetGame.Code)
	}


	go func () {
		err := targetGame.HandleMove(moveMessage);
		if	err != nil {
			fmt.Printf("error handling move: %v\n", err)
		}
	}()
	return nil
}

func (l *Lobby) HandleCreateGame (message []byte, client *socket.Connection) error {
	var createGame GameMessage
	if err := json.Unmarshal(message, &createGame); err != nil {
		return fmt.Errorf("error unmarshalling chat message: %v\n", err)
	}

	newGame := game.NewGame()
	l.CleanupGames(client, newGame.Code)

	l.Games[newGame.Code] = newGame

	newGame.Players[client.Id] = client

	fmt.Printf("total games: %v\n", len(l.Games))

	msg := &GameMessage{
		Type: "game-code",
		Code: newGame.Code,
	}
	client.Send <- msg

	return nil
}

func (l *Lobby) HandleRematch (message []byte) error {
	var rematchMessage GameMessage
	if err := json.Unmarshal(message, &rematchMessage); err != nil {
		return fmt.Errorf("error unmarshalling move-message: %v\n", err)
	}

	targetGame, exists := l.Games[rematchMessage.Code]
	if !exists {
		return fmt.Errorf("game doesn't exist at: %v", rematchMessage.Code)
	}

	targetGame.RematchPool = targetGame.RematchPool + 1

	if targetGame.RematchPool == 2 {
		targetGame.Reset()
		targetGame.Start(true)
	} else {
		sendRematchRequest(targetGame.Players)
	}
	return nil
}

func sendRematchRequest(players map[string]*socket.Connection) {
	for _, player := range players {
		rematchReq := GameMessage{
			Type: "rematch-request",
		}
		player.Send <- rematchReq
	}
}
