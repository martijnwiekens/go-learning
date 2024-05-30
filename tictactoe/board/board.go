package board

import "fmt"

type Board struct {
	board [][]uint8
	size  uint8
}

func NewBoard(boardSize uint8) *Board {
	// Generate board
	board := make([][]uint8, boardSize)
	for i := range board {
		board[i] = make([]uint8, boardSize)
	}

	return &Board{
		board: board,
		size:  boardSize,
	}
}

func (playBoard *Board) GetBoard() [][]uint8 {
	return playBoard.board
}

func (playBoard *Board) SetPosition(row uint8, col uint8, value uint8) {
	playBoard.board[row][col] = value
}

func (playBoard *Board) GetPosition(row uint8, col uint8) uint8 {
	return playBoard.board[row][col]
}

func (playBoard *Board) GetBoardSize() int {
	return int(playBoard.size)
}

func (playBoard *Board) IsFull() bool {
	for i := 0; i < int(playBoard.size); i++ {
		for j := 0; j < int(playBoard.size); j++ {
			if playBoard.board[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

func (playBoard *Board) CheckWin(player uint8) bool {
	// Create list of possible winning combinations
	/*
				[
					[[0,0], [0,1], [0,2]], Row #1
					[[1,0], [1,1], [1,2]], Row #2
					[[2,0], [2,1], [2,2]], Row #3
		            [[0,0], [1,0], [2,0]], Col #1
		            [[0,1], [1,1], [2,1]], Col #2
		            [[0,2], [1,2], [2,2]], Col #3
		            [[0,0], [1,1], [2,2]], Diagonal #1
					[[2,0], [1,1], [0,2]], Diagonal #2
				]
	*/
	boardSize := playBoard.GetBoardSize()
	var winningCombinations [][][2]uint8

	// First up: rows
	for i := 0; i < boardSize; i++ {
		rowCombinations := [][2]uint8{}
		for j := 0; j < boardSize; j++ {
			rowCombinations = append(rowCombinations, [2]uint8{uint8(i), uint8(j)})
		}
		winningCombinations = append(winningCombinations, rowCombinations)
	}

	if len(winningCombinations) != boardSize {
		fmt.Println("Error: winningCombinations #1 is not the same length as boardSize")
		fmt.Println(winningCombinations)
		return false
	}

	// Second up: columns
	for i := 0; i < boardSize; i++ {
		rowCombinations := [][2]uint8{}
		for j := 0; j < boardSize; j++ {
			rowCombinations = append(rowCombinations, [2]uint8{uint8(j), uint8(i)})
		}
		winningCombinations = append(winningCombinations, rowCombinations)
	}

	if len(winningCombinations) != boardSize*2 {
		fmt.Println("Error: winningCombinations #2 is not the same length as boardSize")
		fmt.Println(winningCombinations)
		return false
	}

	// Third up: diagonals
	// They need to be dynamic based on the size of the board
	// First forward dialog
	rowCombinations := [][2]uint8{}
	for i := 0; i < boardSize; i++ {
		rowCombinations = append(rowCombinations, [2]uint8{uint8(i), uint8(i)})
	}
	winningCombinations = append(winningCombinations, rowCombinations)

	// Backward diagonal
	rowCombinations2 := [][2]uint8{}
	for i := 0; i < boardSize; i++ {
		rowCombinations2 = append(rowCombinations2, [2]uint8{uint8(i), uint8(boardSize - 1 - i)})
	}
	winningCombinations = append(winningCombinations, rowCombinations2)

	// Check each winning combination
	// The same player should have all places in the a combination

	if len(winningCombinations) != ((boardSize * 2) + 2) {
		fmt.Println("Error: winningCombinations #3 is not the same length as boardSize")
		fmt.Println(winningCombinations)
		return false
	}

	// Check each winning combination
	// The same player should have all places in the a combination
	for _, winningCombination := range winningCombinations {
		countPlaces := [3]int{0, 0, 0}
		for _, place := range winningCombination {
			holder := playBoard.GetPosition(place[0], place[1])
			countPlaces[holder]++
		}
		if countPlaces[player] == boardSize {
			return true
		}
	}
	return false
}
