package human

import (
	"github.com/martijnwiekens/go-learning/tictactoe/board"
	"github.com/martijnwiekens/go-learning/tictactoe/ui"
)

type HumanPlayer struct {
}

func (humanPlayer *HumanPlayer) AskForMove(playBoard *board.Board) (uint8, uint8) {
	return ui.AskMove()
}
