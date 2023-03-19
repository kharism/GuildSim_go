package main

import (
	"github/kharism/GuildSim_go/internal/cards"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type EbitenCard struct {
	image     *ebiten.Image
	card      cards.Card
	oriWidth  int
	oriHeight int
	// current position
	x int
	y int
	// velocity of card movement
	vx int
	vy int
	// target position if card moved
	tx int
	ty int
}

func (e *EbitenCard) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(0, 0)
	// op.GeoM.Translate(MAIN_DECK_X, MAIN_DECK_Y)
	op.GeoM.Translate(float64(e.x), float64(e.y))
	// op.GeoM.Translate(float64(e.x), float64(e.y))
	screen.DrawImage(e.image, op)
}
func (e *EbitenCard) Update() {
	e.x += e.vx
	e.y += e.vy
	// fmt.Println(e.x, e.y)
	if math.Abs(float64(e.tx-e.x))+math.Abs(float64(e.ty-e.y)) < 15 {
		e.x = e.tx
		e.y = e.ty
		e.vx = 0
		e.vy = 0
	}
}
