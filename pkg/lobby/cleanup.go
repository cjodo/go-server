package lobby

import (
	"fmt"

	"github.com/cjodo/go-server/pkg/socket"
)

// helper for create game
func (l *Lobby) CleanupGames(client *socket.Connection, newCode string) {
	for gameId, game := range l.Games {

		if gameId == newCode {
			fmt.Println("new code")
			return
		}

		if _, exists := game.Players[client.Id]; exists {
			fmt.Println("existing game found")
			delete(l.Games, gameId)
		}
	}
}
