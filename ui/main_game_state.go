package main

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/gamestate"
	"log"
	"math"

	csg "github.com/kharism/golang-csg/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MainGameState struct {
	bgImage          *ebiten.Image
	cardInHand       []*EbitenCard
	stateChanger     AbstractStateChanger
	defaultGamestate *gamestate.DefaultGamestate
}
type OnDrawAction struct {
	mainGameState *MainGameState
}

func (d *OnDrawAction) DoAction(data map[string]interface{}) {
	fmt.Println("OnDrawAction")
	drawnCards := data[cards.EVENT_ATTR_CARD_DRAWN].(cards.Card)
	newEbitenCard := NewEbitenCardFromCard(drawnCards)
	ll := mainGame.(*MainGameState)
	indexCard := len(ll.defaultGamestate.CardsInHand) - 1
	newEbitenCard.x = int(math.Floor(MAIN_DECK_X))
	newEbitenCard.y = int(math.Floor(MAIN_DECK_Y))
	newEbitenCard.tx = int(math.Floor(HAND_START_X + float64(indexCard)*HAND_DIST_X))
	newEbitenCard.ty = HAND_START_Y
	vx := float64(newEbitenCard.tx - newEbitenCard.x)
	vy := float64(newEbitenCard.ty - newEbitenCard.y)
	speedVector := csg.NewVector(vx, vy, 0)
	speedVector = speedVector.Normalize().MultiplyScalar(5)
	newEbitenCard.vx = int(speedVector.X)
	newEbitenCard.vy = int(speedVector.Y)
	fmt.Println(newEbitenCard.x, newEbitenCard.y, newEbitenCard.tx, newEbitenCard.ty)
	ll.cardInHand = append(ll.cardInHand, newEbitenCard)
}
func NewMainGameState(stateChanger AbstractStateChanger) AbstractEbitenState {
	background, _, err = ebitenutil.NewImageFromFile("img/background.png")
	if err != nil {
		log.Fatal(err)
	}
	// image1, _, err := ebitenutil.NewImageFromFile("img/RookieAdventurer.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	cardInHand := []*EbitenCard{}
	// for i := 0; i < 5; i++ {
	// 	card1 := &EbitenCard{image: image1, oriHeight: 600, oriWidth: 450}
	// 	card1.x = int(math.Floor(MAIN_DECK_X))
	// 	card1.y = int(math.Floor(MAIN_DECK_Y))

	// 	card1.tx = int(math.Floor(HAND_START_X + float64(i)*HAND_DIST_X))
	// 	card1.ty = HAND_START_Y

	// 	//velocity
	// 	card1.vx = (card1.tx - card1.x) / 20
	// 	card1.vy = (card1.ty - card1.y) / 20
	// 	cardInHand = append(cardInHand, card1)
	// }
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
	if len(m.cardInHand) > 0 {
		msg := fmt.Sprintf("Card1Pos=(%d,%d)\nCard1Target=(%d,%d)\nCard1V=(%d,%d)", m.cardInHand[0].x, m.cardInHand[0].y,
			m.cardInHand[0].tx, m.cardInHand[1].ty, m.cardInHand[0].vx, m.cardInHand[1].vy)
		ebitenutil.DebugPrint(screen, msg)
	}

}
func (m *MainGameState) Update() error {
	for _, c := range m.cardInHand {
		c.Update()
	}
	return nil
}
