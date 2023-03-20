package main

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/gamestate"
	"image/color"
	"log"
	"math"

	csg "github.com/kharism/golang-csg/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type MainGameState struct {
	bgImage          *ebiten.Image
	bgImage2         *ebiten.Image
	paperBg          *ebiten.Image
	checkMark        *ebiten.Image
	btn              *ebiten.Image
	cardInHand       []*EbitenCard
	stateChanger     AbstractStateChanger
	detailViewCard   *EbitenCard
	defaultGamestate *gamestate.DefaultGamestate

	// sub-states
	currentSubState SubState
	mainState       *mainMainState
	detailState     *detailState
	cardPicker      *cardPickState
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
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		cardInHand := s.m.defaultGamestate.CardsInHand
		fmt.Println("Space pressed")
		go func() {
			cardPick := s.m.defaultGamestate.GetCardPicker().PickCard(cardInHand, "Card from hand")
			fmt.Println("DDDD", cardPick)
		}()
		s.m.currentSubState = s.m.cardPicker
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

type cardPickState struct {
	m            *MainGameState
	cards        []cards.Card
	selectedCard *EbitenCard

	selectedIndex int
	pickedCards   chan (int)
}

func (c *cardPickState) PickCard(list []cards.Card, message string) int {
	c.cards = list
	fmt.Println("Tunggu hasil")
	pickedCards := <-c.pickedCards
	fmt.Println("Dapat hasil", pickedCards)
	return pickedCards
}
func (c *cardPickState) PickCardOptional(list []cards.Card, message string) int {

	return 0
}

const (
	CARDPICKER_START_X = 160
	CARDPICKER_START_Y = 40
)

func (c *cardPickState) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(c.m.bgImage2, op)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(1.3, 1.0)
	op2.GeoM.Translate(120, 0)
	screen.DrawImage(c.m.paperBg, op2)
	op3 := &ebiten.DrawImageOptions{}
	colPerRow := 6
	cardList := []*EbitenCard{}
	for idx, cc := range c.cards {
		ebitenCard := NewEbitenCardFromCard(cc)
		op3.GeoM.Reset()
		op3.GeoM.Scale(HAND_SCALE, HAND_SCALE)
		col := (idx % colPerRow) + 1
		row := (idx / colPerRow) + 1
		// fmt.Println(row, col)
		ebitenCard.x = CARDPICKER_START_X * col
		ebitenCard.y = CARDPICKER_START_Y * row
		op3.GeoM.Translate(float64(CARDPICKER_START_X*col), float64(CARDPICKER_START_Y*row))
		screen.DrawImage(ebitenCard.image, op3)
		cardList = append(cardList, ebitenCard)
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		xCur, yCur := ebiten.CursorPosition()
		fmt.Println("DDDDDD", xCur, yCur)

		for _, ec := range cardList {
			// fmt.Println(ec.x, ec.y)
			if xCur > ec.x && xCur < ec.x+int(math.Floor(ORI_CARD_WIDTH*HAND_SCALE)) &&
				yCur > ec.y && yCur < ec.y+int(math.Floor(ORI_CARD_HEIGHT*HAND_SCALE)) {
				c.selectedCard = ec
				// fmt.Println("Sel", c.selectedCard)
			}
		}
		// check if OK button is clicked
		if xCur > CARDPICKER_START_X && xCur < CARDPICKER_START_X+190 &&
			yCur > 540 && yCur < 540+49 && c.selectedCard != nil {
			fmt.Println("Click OK", len(c.cards))
			for idx, j := range c.cards {
				if j == c.selectedCard.card {
					fmt.Println("Send stuff", idx)
					c.pickedCards <- idx
					c.m.currentSubState = c.m.mainState
					c.selectedCard = nil

					//close(c.pickedCards)
					break

				}
			}
		}
	}

	if c.selectedCard != nil {
		op3.GeoM.Reset()
		op3.GeoM.Translate(CARDPICKER_START_X, 540)
		screen.DrawImage(c.m.btn, op3)
		text.Draw(screen, "OK", mplusNormalFont, CARDPICKER_START_X+70, 570, color.White)
		op3.GeoM.Reset()
		// op3.GeoM.Scale(4, 4)
		op3.GeoM.Translate(float64(c.selectedCard.x), float64(c.selectedCard.y))
		screen.DrawImage(c.m.checkMark, op3)
	}
	// fmt.Println("===")
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
	paperBg, _, err := ebitenutil.NewImageFromFile("img/misc/paper-plain.png")
	if err != nil {
		log.Fatal(err)
	}
	checkmark, _, err := ebitenutil.NewImageFromFile("img/misc/blue_checkmark.png")
	if err != nil {
		log.Fatal(err)
	}
	btn, _, err := ebitenutil.NewImageFromFile("img/misc/green_button00.png")
	if err != nil {
		log.Fatal(err)
	}
	// image1, _, err := ebitenutil.NewImageFromFile("img/RookieAdventurer.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	cardInHand := []*EbitenCard{}
	mgs := &MainGameState{bgImage2: background2, bgImage: background, cardInHand: cardInHand, stateChanger: stateChanger,
		paperBg: paperBg, checkMark: checkmark, btn: btn}
	mainState := &mainMainState{m: mgs}
	detailState := &detailState{m: mgs}
	cardpicker := &cardPickState{m: mgs, pickedCards: make(chan int)}
	mgs.currentSubState = mainState
	mgs.mainState = mainState
	mgs.detailState = detailState
	mgs.cardPicker = cardpicker
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
