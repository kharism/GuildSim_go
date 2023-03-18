package text

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"os"
)

type exitAction struct{}

func (e *exitAction) DoAction() {
	fmt.Println("Game over")
	os.Exit(0)
}
func AttachGameOverListener(state cards.AbstractGamestate) cards.AbstractGamestate {
	quit := exitAction{}
	gameoverlistener := cards.NewStillAliveListener(state, &quit)
	state.AttachListener(cards.EVENT_TAKE_DAMAGE, gameoverlistener)
	return state
}
