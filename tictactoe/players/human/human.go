package human

import (
	"github.com/martijnwiekens/gotictactoe/board"
	"github.com/martijnwiekens/gotictactoe/ui"
)

type HumanPlayer struct {
}

func (humanPlayer *HumanPlayer) AskForMove(playBoard *board.Board) (uint8, uint8) {
	return ui.AskMove()
}
