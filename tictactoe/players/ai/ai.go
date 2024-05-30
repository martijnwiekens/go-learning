package ai

import (
	"math/rand/v2"

	"github.com/martijnwiekens/gotictactoe/board"
)

type AIPlayer struct {
	Mode string
}

// const AI_PLAYER_MODE string = "RANDOM"
const AI_PLAYER_MODE string = "MIN_MAX"

const (
	EMPTY    = 0
	PLAYER_X = 1
	PLAYER_O = 2
)

func (aiPlayer *AIPlayer) AskForMove(playBoard *board.Board) (uint8, uint8) {
	switch aiPlayer.Mode {
	case "RANDOM":
		return aiPlayer.getRandomMove(playBoard)
	case "MIN_MAX":
		return aiPlayer.getMinMaxMove(playBoard)
	default:
		return aiPlayer.getRandomMove(playBoard)
	}
}

func (aiPlayer *AIPlayer) getRandomMove(playBoard *board.Board) (uint8, uint8) {
	// Find a valid move
	var validMove bool = false
	for !validMove {
		// Generate a random row and column
		newRow := rand.IntN(playBoard.GetBoardSize())
		newCol := rand.IntN(playBoard.GetBoardSize())

		// Check if valid
		if playBoard.GetPosition(uint8(newRow), uint8(newCol)) == 0 {
			validMove = true
			return uint8(newRow), uint8(newCol)
		}
	}
	return 0, 0
}

func (aiPlayer *AIPlayer) getMinMaxMove(playBoard *board.Board) (uint8, uint8) {
	bestScore := -1000
	var move [2]uint8
	for i := 0; i < playBoard.GetBoardSize(); i++ {
		for j := 0; j < playBoard.GetBoardSize(); j++ {
			if playBoard.GetPosition(uint8(i), uint8(j)) == EMPTY {
				playBoard.SetPosition(uint8(i), uint8(j), PLAYER_O)
				score := Minimax(playBoard, 0, false)
				playBoard.SetPosition(uint8(i), uint8(j), EMPTY)
				if score > bestScore {
					bestScore = score
					move[0] = uint8(i)
					move[1] = uint8(j)
				}
			}
		}
	}
	return move[0], move[1]
}

func Minimax(b *board.Board, depth int, isMaximizing bool) int {
	// Check if player 1 has won
	if b.CheckWin(PLAYER_X) {
		return -10 + depth
	}

	// Check if player 2 has won
	if b.CheckWin(PLAYER_O) {
		return 10 - depth
	}

	// Check if there are places left
	if b.IsFull() {
		return 0
	}

	// Maximizing player
	if isMaximizing {
		bestScore := -1000
		for i := 0; i < b.GetBoardSize(); i++ {
			for j := 0; j < b.GetBoardSize(); j++ {
				if b.GetPosition(uint8(i), uint8(j)) == EMPTY {
					b.SetPosition(uint8(i), uint8(j), PLAYER_O)
					score := Minimax(b, depth+1, false)
					b.SetPosition(uint8(i), uint8(j), EMPTY)
					if score > bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	} else {
		bestScore := 1000
		for i := 0; i < b.GetBoardSize(); i++ {
			for j := 0; j < 3; j++ {
				if b.GetPosition(uint8(i), uint8(j)) == EMPTY {
					b.SetPosition(uint8(i), uint8(j), PLAYER_X)
					score := Minimax(b, depth+1, true)
					b.SetPosition(uint8(i), uint8(j), EMPTY)
					if score < bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	}
}
