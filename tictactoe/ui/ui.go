package ui

import (
	"fmt"

	"github.com/martijnwiekens/gotictactoe/board"
)

func PrintBoard(playBoardObj *board.Board) {
	/**
	  |_0___1___2__
	0 |	X	X	X
	1 |	X	X	X
	2 |	X	X	X
	*/
	var playBoard = playBoardObj.GetBoard()

	// Create col header line
	fmt.Println()
	fmt.Print("  |__")
	for i := 0; i < len(playBoard); i++ {
		fmt.Print(i)
		fmt.Print("__|__")
	}
	fmt.Println()

	// Print board
	for i := 0; i < len(playBoard); i++ {
		// Create row header line
		fmt.Print(i, " |")
		fmt.Print("  ")
		for j := 0; j < len(playBoard[i]); j++ {
			if playBoard[i][j] == 0 {
				fmt.Print(" ")
			} else if playBoard[i][j] == 1 {
				fmt.Print("X")
			} else if playBoard[i][j] == 2 {
				fmt.Print("O")
			}
			fmt.Print("  |  ")
		}
		fmt.Println()
	}
}

func PrintWinner(winner uint8, totalTurns uint8) {
	fmt.Println()
	if winner == 1 {
		fmt.Printf("X wins in %d turns!\n", totalTurns)
	} else if winner == 2 {
		fmt.Printf("O wins in %d turns!\n", totalTurns)
	} else if winner == 3 {
		fmt.Printf("Tie in %d turns!\n", totalTurns)
	}
}

func PrintTurn(turn uint8) {
	fmt.Println()
	if turn == 1 {
		fmt.Println("X's turn")
	} else if turn == 2 {
		fmt.Println("O's turn")
	}
}

func AskMove() (uint8, uint8) {
	var row uint8
	var col uint8
	fmt.Print("Enter row: ")
	fmt.Scan(&row)
	fmt.Print("Enter column: ")
	fmt.Scan(&col)
	fmt.Println("_____________________")
	return row, col
}

func WrongMove() {
	fmt.Println("Wrong move! Try again")
}

func AskForRestart() bool {
	// Ask for press ENTER to restart
	fmt.Println("Press ENTER to restart")
	fmt.Scanln()
	return true
}

func PrintIntGame() {
	fmt.Println("---- TicTactToe ----")
}

func PrintStartGame(totalGames uint8) {
	fmt.Printf("Starting new game #%d\n", totalGames)
}
