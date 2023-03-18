package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MainGameState struct {
	bgImage      *ebiten.Image
	cardInHand   []*EbitenCard
	stateChanger AbstractStateChanger
}

func NewMainGameState(stateChanger AbstractStateChanger) AbstractEbitenState {
	background, _, err = ebitenutil.NewImageFromFile("img/background.png")
	if err != nil {
		log.Fatal(err)
	}
	image1, _, err := ebitenutil.NewImageFromFile("img/RookieAdventurer.png")
	if err != nil {
		log.Fatal(err)
	}
	cardInHand := []*EbitenCard{}
	for i := 0; i < 5; i++ {
		card1 := &EbitenCard{image: image1, oriHeight: 600, oriWidth: 450}
		card1.x = int(math.Floor(MAIN_DECK_X))
		card1.y = int(math.Floor(MAIN_DECK_Y))

		card1.tx = int(math.Floor(HAND_START_X + float64(i)*HAND_DIST_X))
		card1.ty = HAND_START_Y

		//velocity
		card1.vx = (card1.tx - card1.x) / 20
		card1.vy = (card1.ty - card1.y) / 20
		cardInHand = append(cardInHand, card1)
	}
	return &MainGameState{bgImage: background, cardInHand: cardInHand, stateChanger: stateChanger}
}
func (m *MainGameState) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(background, op)
	for _, c := range m.cardInHand {
		c.Draw(screen)
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		fmt.Println("oo")
		ShowCardDetail = true
	}
	if ShowCardDetail {

	}
	msg := fmt.Sprintf("Card1Pos=(%d,%d)\nCard1Target=(%d,%d)\nCard1V=(%d,%d)", m.cardInHand[0].x, m.cardInHand[0].y, m.cardInHand[0].tx, m.cardInHand[1].ty, m.cardInHand[0].vx, m.cardInHand[1].vy)
	ebitenutil.DebugPrint(screen, msg)
}
func (m *MainGameState) Update() error {
	for _, c := range m.cardInHand {
		c.Update()
	}
	return nil
}
