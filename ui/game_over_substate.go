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
	alpha        float64
	alphaMove    float64
}

func (c *gameOverSubstate) Update() error {
	// c.alphaMove = 1
	if c.alpha < 1.0 {
		c.alpha += c.alphaMove
		c.alphaMove += 0.0001
		// fmt.Println(c.alpha, c.alphaMove)
	}

	return nil
}
func (c *gameOverSubstate) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(c.m.bgImage2, op)
	op.ColorM.Scale(1, 1, 1, c.alpha)
	screen.DrawImage(c.m.GameOver, op)

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mainGame.(*MainGameState).stateChanger.ChangeState(STATE_MAIN_MENU)
	}
}
