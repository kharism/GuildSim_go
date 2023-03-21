package main

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/gamestate"
	"image/color"
	"log"
	"math"
	"math/rand"

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
	iconCombat       *ebiten.Image
	iconExplore      *ebiten.Image
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
		s.m.detailState.prevSubState = s
		s.m.currentSubState = s.m.detailState
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		// cardInHand := s.m.defaultGamestate.CardsInHand
		fmt.Println("Space pressed")
		go func() {
			kk := []cards.Card{}
			for i := 0; i < 10; i++ {
				adv := cards.NewRookieAdventurer(s.m.defaultGamestate)
				com := cards.NewRookieCombatant(s.m.defaultGamestate)
				kk = append(kk, &adv, &com)
			}
			rand.Shuffle(20, func(i, j int) { kk[i], kk[j] = kk[j], kk[i] })
			cardPick := s.m.defaultGamestate.GetCardPicker().PickCardOptional(kk, "Card from hand")
			fmt.Println("DDDD", cardPick)
		}()
		s.m.currentSubState = s.m.cardPicker
	}
}

type detailState struct {
	m            *MainGameState
	prevSubState SubState
}

func (s *detailState) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(s.m.bgImage2, op)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(600-ORI_CARD_WIDTH/2, 0)
	screen.DrawImage(s.m.detailViewCard.image, op2)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		s.m.currentSubState = s.prevSubState
		s.prevSubState = nil
	}
}

type cardPickState struct {
	m             *MainGameState
	cards         []cards.Card
	selectedCard  *EbitenCard
	optional      bool
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
	c.cards = list
	c.optional = true
	fmt.Println("Tunggu hasil")
	pickedCards := <-c.pickedCards
	fmt.Println("Dapat hasil", pickedCards)
	return pickedCards
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
	iconExplore, _, err := ebitenutil.NewImageFromFile("img/misc/exploration.png")
	if err != nil {
		log.Fatal(err)
	}
	iconCombat, _, err := ebitenutil.NewImageFromFile("img/misc/combat.png")
	if err != nil {
		log.Fatal(err)
	}
	// image1, _, err := ebitenutil.NewImageFromFile("img/RookieAdventurer.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	cardInHand := []*EbitenCard{}
	mgs := &MainGameState{bgImage2: background2, bgImage: background, cardInHand: cardInHand, stateChanger: stateChanger,
		paperBg: paperBg, checkMark: checkmark, btn: btn, iconCombat: iconCombat, iconExplore: iconExplore}
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
	res := m.defaultGamestate.GetCurrentResource()
	hp := m.defaultGamestate.GetCurrentHP()
	text.Draw(screen, fmt.Sprintf("HP %d", hp), mplusResource, 80, 40, color.RGBA{255, 0, 0, 255})

	op.GeoM.Reset()
	op.GeoM.Scale(0.8, 0.8)
	op.GeoM.Translate(350, 0)
	screen.DrawImage(m.iconCombat, op)
	combat, ok := res.Detail[cards.RESOURCE_NAME_COMBAT]
	if !ok {
		combat = 0
	}
	text.Draw(screen, fmt.Sprintf("%d", combat), mplusResource, 500, 40, color.RGBA{255, 0, 0, 255})

	op.GeoM.Reset()
	op.GeoM.Scale(0.8, 0.8)
	op.GeoM.Translate(540, 0)
	screen.DrawImage(m.iconExplore, op)
	explore, ok := res.Detail[cards.RESOURCE_NAME_EXPLORATION]
	if !ok {
		explore = 0
	}
	text.Draw(screen, fmt.Sprintf("%d", explore), mplusResource, 670, 40, color.RGBA{0, 255, 0, 255})

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
