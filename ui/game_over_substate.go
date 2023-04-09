package main

import (
	"github/kharism/GuildSim_go/internal/cards"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type gameOverSubstate struct {
	m            *MainGameState
	cards        []cards.Card
	selectedCard *EbitenCard
}

func (c gameOverSubstate) Update() error {
	return nil
}
func (c gameOverSubstate) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(c.m.bgImage2, op)
	screen.DrawImage(c.m.GameOver, op)

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mainGame.(*MainGameState).stateChanger.ChangeState(STATE_MAIN_MENU)
	}
}
