package game

import (
	"fmt"

	"github.com/cjodo/go-server/pkg/socket"
	"github.com/cjodo/go-server/pkg/util"
)

//defaults
var (
	defaultBoard 	=		[9]string{ "", "", "", "", "", "", "", "", "" }
	defaultTurn 	=		"X" 
	defaultWinner = 	""
)

type Game struct {
	Code					string
	Players 			map[string]*socket.Connection // token -> connection
	Board 				[9]string
	Turn 					string
	Winner 				string
	RematchPool   int
}

type WinnerMessage struct {
	Type 			string 	`json:"type"`
	Game			string 	`json:"game-code"`
	Player 		string	`json:"player"`
}

type startMessage struct {
	Type			string			`json:"type"`
	Game			string			`json:"game-code"`
	Turn 			string			`json:"turn"`
	Player		string			`json:"player"`
	Message		string			`json:"message"`
	Board			[9]string		`json:"board"`
}

type MoveMessage struct {
	Type			string		`json:"type"`
	Code			string 		`json:"game-code"`
	Player		string 		`json:"player"`
	Move			int				`json:"move"`
	Turn			string		`json:"turn"`
	Board			[9]string	`json:"board"`
}

func NewGame () *Game {
	return &Game{
		Code:			util.RandomStringRunes(8),
		Players: 	make(map[string]*socket.Connection),
		Board: 		defaultBoard,
		Turn: 		defaultTurn,
		Winner: 	defaultWinner,
		RematchPool: 0,
	}
}

//Todo: make a notify players method on game: func (g *Game) notifyPlayers(msg interface{}) error {}

func (g *Game) Start(rematch bool) {
	var playerOne string
	var playerTwo string

	if rematch {
		playerOne = "O"
		playerTwo = "X"
	} else {
		playerOne = "X"
		playerTwo = "O"
	}

	currPlayer := playerOne
	for _, player := range g.Players {
		startMessage := &startMessage{
			Type:			"start-game",	
			Game:			g.Code,
			Turn:	 		"X",
			Player: 	currPlayer,
			Message:	"Game has started",
			Board:		[9]string{ "", "", "", "", "", "", "", "", "" },	
		}

		currPlayer = playerTwo

		g.notifyPlayer(startMessage, player)
	}
}

func (g *Game) checkValidMove (move int) bool {
	if g.Board[move] != "" && g.Winner == "" {
		return false
	}

	return true
}

func (g *Game) HandleMove (m MoveMessage) error {
	fmt.Println("message from", m.Turn)

	if g.Turn != m.Turn {
		return fmt.Errorf("not that players turn");
	}

	if valid := g.checkValidMove(m.Move); !valid {
		return fmt.Errorf("not a valid move");
	}

	winner := g.CheckWin(m.Player, m.Move)
	if winner != ""{
		g.Winner = winner
		g.handleWinner(m.Player)
	}

	g.Board[m.Move] = m.Turn
	m.Board = g.Board

	for _, client := range g.Players{
		g.notifyPlayer(m, client)
	}

	g.changeTurn()

	return nil
}

func (g *Game) notifyPlayer (m interface{}, client *socket.Connection){
	fmt.Println("notify: ", client)
	go func() {
		client.Send <- m
	}()
}

func (g *Game) changeTurn() {
	if g.Turn == "X" {
		g.Turn = "O"
	} else {
		g.Turn = "X"
	}
}

func (g *Game) Reset() {
	g.Board = 			defaultBoard
	g.Turn = 				defaultTurn
	g.Winner = 			defaultWinner
	g.RematchPool = 0
}

func (g *Game) handleWinner(player string) {
	var winMessaage = WinnerMessage{
		Type: "winner",
		Game: g.Code,
		Player: player,
	}

	for _, player := range g.Players {
		player.Send <- winMessaage
	}
}

