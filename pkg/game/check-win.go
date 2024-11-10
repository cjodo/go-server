package game

type Vec2 struct {
	X int
	Y int
}

func (g *Game) CheckWin(gameBoard [9]string, player string, move int) string {
	newBoard := buildBoard2D(gameBoard, player, move) 
	//rows
	for row := range 3 {
		left := newBoard[row][0]
		middle := newBoard[row][1]
		right := newBoard[row][2]

		if (left != "" || middle != "" || right != "") {
			if(left == middle && middle == right) {
				return player
			}
		}
	}
	//columns
	for col := range 3 {
		left := newBoard[0][col]
		middle := newBoard[1][col]
		right := newBoard[2][col]

		if (left != "" || middle != "" || right != "") {
			if(left == middle && middle == right) {
				return player
			}
		}
	}

	//diagonals
	UL := newBoard[0][0];
	M := newBoard[1][1];
	BR := newBoard[2][2];
  // Check the first diagonal
	if (UL != "" || M != "" || BR != "") {
    if (UL == M && M == BR) {
      return player;
    }
  }

	UR := newBoard[0][2];
	BL := newBoard[2][0];
	if (UR  != "" || M != "" || BL != "") {
    if (UR == M && M == BL) {
      return player;
    }
  }

	return ""
}

func buildBoard2D (gameBoard [9]string, player string, move int) [][]string {
	var newBoard = make([][]string, 3)

	var playerCoord = Vec2{
		X: move / 3,
		Y: move % 3,
	}

	for j := range 3 {
		newBoard[j] = make([]string, 3)
	}

	for i := range 9 {
		var board = Vec2{
			X: i / 3,
			Y: i % 3,
		}

		if(board.X == playerCoord.X && board.Y == playerCoord.Y) {
			newBoard[board.X][board.Y] = player
		} else {
			newBoard[board.X][board.Y] = gameBoard[i]
		}
	}
	return newBoard
}
