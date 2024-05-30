package human

import (
	"github.com/martijnwiekens/gofourinarow/board"
	"github.com/martijnwiekens/gofourinarow/ui"
)

type HumanPlayer struct {
}

func (humanPlayer *HumanPlayer) AskForMove(playBoard *board.Board) uint8 {
	return ui.AskMove()
}
