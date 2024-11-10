package lobby

import (
	"fmt"

	"github.com/cjodo/go-server/pkg/game"
	"github.com/cjodo/go-server/pkg/socket"
)

func (l *Lobby) handleReconnect(newConn *socket.Connection) {
	l.clients[newConn.Id].Conn = newConn.Conn

	inGame := l.findClientInGame(newConn)

	if inGame != nil {
		started := len(inGame.Players) > 1
		msg := ReconnectMessage{
			Type: "reconnect",
			Game: inGame.Code,
			Board: inGame.Board,
			Winner: inGame.Winner,
			IsStarted: started,
		}

		newConn.Send <- msg
	}

	return
}

func (l *Lobby) findClientInGame(c *socket.Connection) *game.Game {

	for _, g := range l.Games {
		var mostRecent *game.Game
		if _, exists := g.Players[c.Id]; exists {
			if mostRecent == nil {
				fmt.Println("found game")
				mostRecent = g
			}
		}
		return mostRecent
	}

	return nil
}

