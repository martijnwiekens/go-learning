package game

import (
	"github.com/martijnwiekens/go-learning/tictactoe/board"
	"github.com/martijnwiekens/go-learning/tictactoe/players/ai"
	"github.com/martijnwiekens/go-learning/tictactoe/players/human"
	"github.com/martijnwiekens/go-learning/tictactoe/ui"
)

var totalGames uint8 = 0

type Player interface {
	AskForMove(playBoard *board.Board) (uint8, uint8)
}

type Game struct {
	currentPlayerIndicator uint8
	currentPlayer          Player
	totalTurns             uint8
	player1                Player
	player2                Player
	playBoard              *board.Board
}

func StartGame(boardSize uint8, withAi bool) {
	// Print a message
	ui.PrintIntGame()

	// Create the game
	game := Game{
		currentPlayerIndicator: 1,
		currentPlayer:          nil,
		totalTurns:             1,
		player1:                &human.HumanPlayer{},
		player2:                &human.HumanPlayer{},
		playBoard:              board.NewBoard(boardSize),
	}

	// Count the game
	totalGames++

	// Check if we need an AI
	if withAi {
		// Create a new AI
		game.player2 = &ai.AIPlayer{Mode: "MIN_MAX"}
	}

	// Set the current player
	game.currentPlayer = game.player1

	// Start the game loop
	restart := gameLoop(&game)
	if restart {
		ui.PrintStartGame(totalGames)
		StartGame(boardSize, withAi)
	}
}

func gameLoop(gameObj *Game) bool {
	// Find the board
	playBoardObj := gameObj.playBoard

	// Print the board
	ui.PrintBoard(playBoardObj)

	// Print the current player
	ui.PrintTurn(gameObj.currentPlayerIndicator)

	// Keep asking for the right move
	var rightMove bool = false
	for !rightMove {
		// Ask for move
		newRow, newCol := gameObj.currentPlayer.AskForMove(playBoardObj)

		// Check if the move is valid
		rightMove = checkMove(newRow, newCol, gameObj)

		// Wrong move
		if !rightMove {
			ui.WrongMove()
		}
	}

	// Check for winner
	winner := checkWinner(gameObj)
	if winner != 3 {
		// Somebody won the game
		ui.PrintWinner(winner, gameObj.totalTurns)

		// Print the winning board
		ui.PrintBoard(playBoardObj)

		// Ask for restart
		return ui.AskForRestart()
	}

	// Switch player
	changePlayer(gameObj)

	// Increase turn
	gameObj.totalTurns++

	// Start the loop again
	return gameLoop(gameObj)
}

func checkMove(newRow uint8, newCol uint8, gameObj *Game) bool {
	// Get the board
	playBoardObj := gameObj.playBoard

	// Check if the newRow is valid
	if newRow >= uint8(playBoardObj.GetBoardSize()) {
		return false
	}

	// Check if the newCol is valid
	if newCol >= uint8(playBoardObj.GetBoardSize()) {
		return false
	}

	// Check if the move is valid
	if playBoardObj.GetPosition(newRow, newCol) == 0 {
		playBoardObj.SetPosition(newRow, newCol, gameObj.currentPlayerIndicator)
		return true
	} else {
		return false
	}
}

func checkWinner(gameObj *Game) uint8 {
	// Get the board
	playBoardObj := gameObj.playBoard

	// Check the winner
	if playBoardObj.CheckWin(gameObj.currentPlayerIndicator) {
		return gameObj.currentPlayerIndicator
	}

	// Check for tie
	if playBoardObj.IsFull() {
		return 0
	}

	return 3
}

func changePlayer(gameObj *Game) {
	if gameObj.currentPlayerIndicator == 1 {
		gameObj.currentPlayerIndicator = 2
		gameObj.currentPlayer = gameObj.player2
	} else if gameObj.currentPlayerIndicator == 2 {
		gameObj.currentPlayerIndicator = 1
		gameObj.currentPlayer = gameObj.player1
	}
}
