package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type actClearSubstate struct {
	m *MainGameState
	// cards        []cards.Card
	// selectedCard *EbitenCard
	alpha     float64
	alphaMove float64
	doneFunc  func()
}

func (c *actClearSubstate) Update() error {
	// c.alphaMove = 1
	if c.alpha < 1.0 {
		c.alpha += c.alphaMove
		c.alphaMove += 0.0001
		// fmt.Println(c.alpha, c.alphaMove)
	}

	return nil
}
func (c *actClearSubstate) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(c.m.bgImage2, op)
	op.ColorM.Scale(1, 1, 1, c.alpha)
	screen.DrawImage(c.m.ActClear, op)

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		//mainGame.(*MainGameState).stateChanger.ChangeState(STATE_MAIN_GAME)
		go c.doneFunc()

	}
}
