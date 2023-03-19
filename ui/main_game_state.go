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
	bgImage2         *ebiten.Image
	cardInHand       []*EbitenCard
	stateChanger     AbstractStateChanger
	detailViewCard   *EbitenCard
	defaultGamestate *gamestate.DefaultGamestate
	currentSubState  SubState
	mainState        *mainMainState
	detailState      *detailState
}
type SubState interface {
	Draw(screen *ebiten.Image)
}
type mainMainState struct {
	m *MainGameState
}

func (s *mainMainState) Draw(screen *ebiten.Image) {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		xCur, yCur := ebiten.CursorPosition()
		fmt.Println("oo", xCur, yCur)
		if yCur > HAND_START_Y {
			// right click on hand
			for i := len(s.m.cardInHand) - 1; i >= 0; i-- {
				if s.m.cardInHand[i].x < xCur {
					s.m.detailViewCard = s.m.cardInHand[i]
					fmt.Println("cardIndex at", i)
					break
				}
			}
		}
		s.m.currentSubState = s.m.detailState
	}
}

type detailState struct {
	m *MainGameState
}

func (s *detailState) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(s.m.bgImage2, op)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(600-ORI_CARD_WIDTH/2, 0)
	screen.DrawImage(s.m.detailViewCard.image, op2)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		s.m.currentSubState = s.m.mainState
	}
}

type OnDrawAction struct {
	mainGameState *MainGameState
}

const (
	CARD_MOVE_SPEED = 5
)

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
	speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
	newEbitenCard.vx = int(speedVector.X)
	newEbitenCard.vy = int(speedVector.Y)
	fmt.Println(newEbitenCard.x, newEbitenCard.y, newEbitenCard.tx, newEbitenCard.ty)
	ll.cardInHand = append(ll.cardInHand, newEbitenCard)
}
func NewMainGameState(stateChanger AbstractStateChanger) AbstractEbitenState {
	background, _, err := ebitenutil.NewImageFromFile("img/background.png")
	if err != nil {
		log.Fatal(err)
	}
	background2, _, err := ebitenutil.NewImageFromFile("img/background_trans.png")
	if err != nil {
		log.Fatal(err)
	}
	// image1, _, err := ebitenutil.NewImageFromFile("img/RookieAdventurer.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	cardInHand := []*EbitenCard{}
	mgs := &MainGameState{bgImage2: background2, bgImage: background, cardInHand: cardInHand, stateChanger: stateChanger}
	mainState := &mainMainState{m: mgs}
	detailState := &detailState{m: mgs}
	mgs.currentSubState = mainState
	mgs.mainState = mainState
	mgs.detailState = detailState
	return mgs
}

var ShowCardDetail = false
var ShowCardPicker = false

func (m *MainGameState) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(m.bgImage, op)
	for _, c := range m.cardInHand {
		c.Draw(screen)
	}

	m.currentSubState.Draw(screen)

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
