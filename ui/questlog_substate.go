package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type questLogState struct {
	m *MainGameState
	// cards         []cards.Card
	// selectedCard  *EbitenCard
	// optional      bool
	// selectedIndex int
	// pickedCards   chan (int)
	message string
}

func (c *questLogState) Update() error {
	return nil
}
func (c *questLogState) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(c.m.bgImage2, op)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(1.3, 1.0)
	op2.GeoM.Translate(120, 0)
	screen.DrawImage(c.m.paperBg, op2)
	text.Draw(screen, c.message, mplusNormalFont, CARDPICKER_START_X, 40, color.RGBA{255, 255, 255, 255})
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		c.m.currentSubState = c.m.mainState
		// xCurInt, yCurInt := ebiten.CursorPosition()
		// //fmt.Println("DDDDDD", xCur, yCur)
		// xCur, yCur := float64(xCurInt), float64(yCurInt)
		// if xCur > CARDPICKER_START_X+200 && xCur < CARDPICKER_START_X+200+190 &&
		// 	yCur > 540 && yCur < 540+49 {
		// 	// c.pickedCards <- -1
		// 	c.m.currentSubState = c.m.mainState
		// 	// c.selectedCard = nil
		// }
	}
}
