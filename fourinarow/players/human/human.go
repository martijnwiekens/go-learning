package human

import (
	"github.com/martijnwiekens/go-learning/fourinarow/board"
	"github.com/martijnwiekens/go-learning/fourinarow/ui"
)

type HumanPlayer struct {
}

func (humanPlayer *HumanPlayer) AskForMove(playBoard *board.Board) uint8 {
	return ui.AskMove()
}
