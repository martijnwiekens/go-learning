package ui

import (
	"fmt"

	"github.com/martijnwiekens/gofourinarow/board"
)

func PrintBoard(playBoardObj *board.Board) {
	var playBoard = playBoardObj.GetBoard()

	// Print board
	fmt.Println()
	for i := 0; i < len(playBoard); i++ {
		// Create row header line
		fmt.Print("|  ")
		for j := 0; j < len(playBoard[i]); j++ {
			if playBoard[j][i] == 0 {
				fmt.Print(" ")
			} else if playBoard[j][i] == 1 {
				fmt.Print("X")
			} else if playBoard[j][i] == 2 {
				fmt.Print("O")
			}
			fmt.Print("  |  ")
		}
		fmt.Println()
	}

	// Create col header line
	fmt.Println("|_____|_____|_____|_____")
	for i := 0; i < len(playBoard); i++ {
		fmt.Print("|  ", i, "  ")
	}
	fmt.Println()
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

func AskMove() uint8 {
	var row uint8
	fmt.Print("Enter row: ")
	fmt.Scan(&row)
	fmt.Println("_____________________")
	return row
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
