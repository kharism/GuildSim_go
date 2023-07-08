package main

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/gamestate"
	"image/color"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	csg "github.com/kharism/golang-csg/core"
	"golang.org/x/image/font"
)

type MainGameState struct {
	bgImage       *ebiten.Image
	bgImage2      *ebiten.Image
	paperBg       *ebiten.Image
	checkMark     *ebiten.Image
	btn           *ebiten.Image
	iconCombat    *ebiten.Image
	iconExplore   *ebiten.Image
	DiscardPile   *ebiten.Image
	MainDeck      *ebiten.Image
	EndturnBtn    *ebiten.Image
	GameOver      *ebiten.Image
	ItemIcon      *ebiten.Image
	Reputation    *ebiten.Image
	Block         *ebiten.Image
	cardsInCenter []*EbitenCard
	cardInHand    []*EbitenCard
	cardsPlayed   []*EbitenCard
	startDragX    int
	startDragY    int
	dragMode      bool
	dragDist      int
	stillAnim     bool
	// cards in limbo meaning cards that is moving into cooldownpile or banished pile
	// they have still visible until they reach those position
	cardsInLimbo     []*EbitenCard
	textInLimbo      []*EbitenText
	stateChanger     AbstractStateChanger
	detailViewCard   *EbitenCard
	mutex            *sync.Mutex
	defaultGamestate *gamestate.DefaultGamestate
	limiter          string

	// ui related stuff so we don't do mutex lock every update/draw
	hp              int
	combat          int
	exploration     int
	block           int
	reputation      int
	NumCardInDeck   int
	NumCardCooldown int

	// channels
	CardPlayedChan chan cards.Card

	// sub-states
	currentSubState SubState
	mainState       *mainMainState
	detailState     *detailState
	cardPicker      *cardPickState
	boolPicker      *boolPickState
	cardListState   *cardListState
	gameoverState   *gameOverSubstate
}
type SubState interface {
	Draw(screen *ebiten.Image)
	Update() error
}
type mainMainState struct {
	m     *MainGameState
	mutex *sync.Mutex // mutex to allows some order of animation
}

func CreateDoneFunc(c *EbitenCard, wg *sync.WaitGroup) func() {
	return func() {
		//fmt.Println("Sending done signal", c.card.GetName())
		wg.Done()
	}
}

// this routine play cards, any cards played will be sent to this routine
func CardPlayer(m *MainGameState, cardPlayed <-chan cards.Card) {
	for card := range cardPlayed {
		m.defaultGamestate.PlayCard(card)
	}
}
func (s *mainMainState) Reset() {
	s.m.cardInHand = []*EbitenCard{}
	s.m.cardsInLimbo = []*EbitenCard{}
	s.m.cardsInCenter = []*EbitenCard{}
}
func (s *mainMainState) Draw(screen *ebiten.Image) {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		xCurInt, yCurInt := ebiten.CursorPosition()
		// fmt.Println("oo", xCur, yCur)
		xCur, yCur := float64(xCurInt), float64(yCurInt)
		cardCollection := []*EbitenCard{}
		if yCur > HAND_START_Y {
			// right click on hand
			cardCollection = s.m.cardInHand
		} else if yCur > PLAYED_START_Y {
			cardCollection = s.m.cardsPlayed
		} else if yCur > CENTER_DECK_START_Y {
			cardCollection = s.m.cardsInCenter
		}
		for i := len(cardCollection) - 1; i >= 0; i-- {
			if cardCollection[i].x < xCur {
				s.m.detailViewCard = cardCollection[i]
				//fmt.Println("cardIndex at", i)
				break
			}
		}
		if s.m.detailViewCard != nil {
			s.m.detailState.prevSubState = s
			s.m.currentSubState = s.m.detailState
		}

	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.m.dragMode = true
		s.m.startDragX, _ = ebiten.CursorPosition()
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && !s.m.dragMode {
		//fmt.Println(inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft))
		xCurInt, yCurInt := ebiten.CursorPosition()
		xCur, yCur := float64(xCurInt), float64(yCurInt)
		if xCur != float64(s.m.startDragX) {
			// it basically release from drag mode and not picking card
			// if the first/leftmost card in hand is on the right of HAND_START, move all hand card to their ori position

		} else if xCur > DISCARD_START_X && xCur < DISCARD_START_X+HAND_SCALE*ORI_CARD_WIDTH {
			s.m.cardListState.cards = s.m.defaultGamestate.CardsDiscarded.List()
			s.m.currentSubState = s.m.cardListState
		} else if xCur > ENDTURN_START_X {
			fmt.Println("Endturn")
			if s.m.stillAnim {
				return
			} else {
				s.m.stillAnim = true
				go func() {
					s.m.defaultGamestate.EndTurn()

					s.m.defaultGamestate.BeginTurn()
					s.m.mutex.Lock()
					s.m.stillAnim = false
					s.m.mutex.Unlock()
				}()
			}

		} else if yCur > HAND_START_Y && xCur < DISCARD_START_X && xCur >= HAND_START_X {
			// left click on hand
			for i := len(s.m.cardInHand) - 1; i >= 0; i-- {
				if s.m.cardInHand[i].x < xCur {
					//s.m.defaultGamestate.PlayCard(s.m.cardInHand[i].card)
					s.m.CardPlayedChan <- s.m.cardInHand[i].card
					//s.m.detailViewCard = s.m.cardInHand[i]
					//fmt.Println("cardIndex at", i)
					break
				}
			}
		} else if yCur > HAND_START_Y && xCur < HAND_START_X {
			fmt.Println("Clicked deck")
			// left click on main deck, look at the content of main deck
			s.m.cardListState.cards = []cards.Card{}
			cardInDeck := s.m.defaultGamestate.CardsInDeck.List()
			for _, v := range cardInDeck {
				s.m.cardListState.cards = append(s.m.cardListState.cards, v)
			}
			rand.Shuffle(len(s.m.cardListState.cards), func(i, j int) {
				s.m.cardListState.cards[i], s.m.cardListState.cards[j] = s.m.cardListState.cards[j], s.m.cardListState.cards[i]
			})
			s.m.currentSubState = s.m.cardListState
		} else if yCur > PLAYED_START_Y {
			// do nothing. This is just so we have safe area to release left mouse button
		} else if yCur > CENTER_START_Y {
			// recruite/explore/defeat card from center row
			for i := len(s.m.cardsInCenter) - 1; i >= 0; i-- {
				if s.m.cardsInCenter[i].x < xCur {
					clickedCard := s.m.cardsInCenter[i]
					switch clickedCard.card.GetCardType() {
					case cards.Area:
						go s.m.defaultGamestate.Explore(clickedCard.card)
					case cards.Hero:
						go s.m.defaultGamestate.RecruitCard(clickedCard.card)
					case cards.Monster:
						if _, ok := clickedCard.card.(cards.Recruitable); ok {
							go func() {
								recruit := s.m.defaultGamestate.GetBoolPicker().BoolPick("Recruit " + clickedCard.card.GetName() + " ?")
								if recruit {
									s.m.defaultGamestate.RecruitCard((clickedCard.card))
								} else {
									s.m.defaultGamestate.DefeatCard(clickedCard.card)
								}
							}()
						} else {
							fmt.Println("Unrecruitable")
							go func() {
								s.m.defaultGamestate.DefeatCard(clickedCard.card)
							}()
						}

					case cards.Trap:
						go s.m.defaultGamestate.Disarm(clickedCard.card)
					}
					//go s.m.defaultGamestate.PlayCard(s.m.cardInHand[i].card)
					//s.m.detailViewCard = s.m.cardInHand[i]
					//fmt.Println("cardIndex at", i)
					break
				}
			}
		} else {
			if xCur > ITEM_ICON_START_X {
				go func() {
					items := s.m.defaultGamestate.ItemCards
					pickedIndex := s.m.cardPicker.PickCardOptional(items, "Items")
					if pickedIndex > -1 {
						item := s.m.defaultGamestate.ItemCards[pickedIndex]
						if _, ok := item.(cards.Consumable); ok {
							s.m.defaultGamestate.RemoveItemIndex(pickedIndex)
							s.m.defaultGamestate.ConsumeItem(item.(cards.Consumable))
						}
					}
				}()
			}
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		// cardInHand := s.m.defaultGamestate.CardsInHand
		fmt.Println("Space pressed")
		// go func() {
		// 	kk := []cards.Card{}
		// 	for i := 0; i < 10; i++ {
		// 		adv := cards.NewRookieAdventurer(s.m.defaultGamestate)
		// 		com := cards.NewRookieCombatant(s.m.defaultGamestate)
		// 		kk = append(kk, &adv, &com)
		// 	}
		// 	rand.Shuffle(20, func(i, j int) { kk[i], kk[j] = kk[j], kk[i] })
		// 	cardPick := s.m.defaultGamestate.GetCardPicker().PickCardOptional(kk, "Card from hand")
		// 	fmt.Println("DDDD", cardPick)
		// }()
		// s.m.currentSubState = s.m.boolPicker
		go func() {
			if s.m.boolPicker.BoolPick("Draw a card?") {
				s.m.defaultGamestate.Draw()
			}
		}()

	}
}

func (s *mainMainState) Update() error {
	for _, c := range s.m.cardInHand {
		c.Update()
	}
	for _, c := range s.m.cardsPlayed {
		c.Update()
	}
	for _, c := range s.m.cardsInCenter {
		c.Update()
	}
	newCardInLimbo := []*EbitenCard{}
	for _, c := range s.m.cardsInLimbo {
		c.Update()
		if c.tx != c.x || c.ty != c.y {
			newCardInLimbo = append(newCardInLimbo, c)
		}
	}
	s.m.cardsInLimbo = newCardInLimbo
	return nil
}

type detailState struct {
	m            *MainGameState
	prevSubState SubState
	wait         chan bool
}

func (s *detailState) ShowDetail(c cards.Card) {
	s.m.detailViewCard = NewEbitenCardFromCard(c)
	s.m.mutex.Lock()
	s.prevSubState = s.m.currentSubState

	s.m.currentSubState = s
	s.m.mutex.Unlock()
	// <-s.wait
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
		// s.wait <- true
		s.prevSubState = nil
	}
}
func (s *detailState) Update() error {
	return nil
}

const (
	CARD_MOVE_SPEED = 10
)

type onLimiterAttach struct {
	mainGameState *MainGameState
}

func (p *onLimiterAttach) DoAction(data map[string]interface{}) {
	fmt.Println("Attach Limiter")
	limiter := data[cards.EVENT_ATTR_LIMITER].(cards.LegalChecker)
	if _, ok := limiter.(fmt.Stringer); ok {
		j := limiter.(fmt.Stringer)
		p.mainGameState.limiter = j.String()
	} else {
		fmt.Println("NOK")
	}
}

type onLimiterDetach struct {
	mainGameState *MainGameState
}

func (p *onLimiterDetach) DoAction(data map[string]interface{}) {
	fmt.Println("Detach Limiter")
	limiter := data[cards.EVENT_ATTR_LIMITER].(cards.LegalChecker)
	if _, ok := limiter.(fmt.Stringer); ok {
		// j := limiter.(fmt.Stringer)
		p.mainGameState.limiter = ""
	} else {
		fmt.Println("NOK")
	}
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
	iconReputation, _, err := ebitenutil.NewImageFromFile("img/misc/reputation.png")
	if err != nil {
		log.Fatal(err)
	}
	iconBlock, _, err := ebitenutil.NewImageFromFile("img/misc/shield.png")
	if err != nil {
		log.Fatal(err)
	}
	mainDeck, _, err := ebitenutil.NewImageFromFile("img/misc/main_deck.png")
	if err != nil {
		log.Fatal(err)
	}
	discardPile, _, err := ebitenutil.NewImageFromFile("img/misc/cool_down.png")
	if err != nil {
		log.Fatal(err)
	}
	EndturnBtn, _, err := ebitenutil.NewImageFromFile("img/misc/end_turn.png")
	if err != nil {
		log.Fatal(err)
	}
	game_over, _, err := ebitenutil.NewImageFromFile("img/game_over.png")
	if err != nil {
		log.Fatal(err)
	}
	item_icon, _, err := ebitenutil.NewImageFromFile("img/misc/bag.png")
	if err != nil {
		log.Fatal(err)
	}
	mutex := &sync.Mutex{}
	// image1, _, err := ebitenutil.NewImageFromFile("img/RookieAdventurer.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	cardInHand := []*EbitenCard{}
	cardsPlayed := []*EbitenCard{}
	mgs := &MainGameState{bgImage2: background2, bgImage: background, cardInHand: cardInHand, stateChanger: stateChanger,
		paperBg: paperBg, checkMark: checkmark, btn: btn, iconCombat: iconCombat, iconExplore: iconExplore, mutex: mutex,
		cardsPlayed: cardsPlayed, DiscardPile: discardPile, MainDeck: mainDeck, EndturnBtn: EndturnBtn, GameOver: game_over,
		ItemIcon: item_icon, Reputation: iconReputation, Block: iconBlock,
	}
	mgs.CardPlayedChan = make(chan cards.Card)
	go CardPlayer(mgs, mgs.CardPlayedChan)
	mainState := &mainMainState{m: mgs, mutex: &sync.Mutex{}}
	detailState := &detailState{m: mgs}
	cardpicker := &cardPickState{m: mgs, pickedCards: make(chan int)}
	boolPicker := &boolPickState{m: mgs, pickedOption: make(chan bool)}
	cardListState := &cardListState{m: mgs}
	mgs.gameoverState = &gameOverSubstate{m: mgs}
	mgs.currentSubState = mainState
	mgs.mainState = mainState
	mgs.detailState = detailState
	mgs.cardPicker = cardpicker
	mgs.cardListState = cardListState
	mgs.boolPicker = boolPicker
	return mgs
}

var ShowCardDetail = false
var ShowCardPicker = false

type EbitenText struct {
	text string
	face font.Face
	// current position
	x float64
	y float64
	// velocity of card movement
	vx float64
	vy float64
	// target position if card moved
	tx             float64
	ty             float64
	color          color.RGBA
	CurrMove       *MoveAnimation
	AnimationQueue []*MoveAnimation
}

func (m *EbitenText) Draw(screen *ebiten.Image) {
	text.Draw(screen, m.text, mplusResource, int(m.x), int(m.y), m.color)
}
func (e *EbitenText) Update() {
	if e.CurrMove == nil && len(e.AnimationQueue) > 0 {
		e.CurrMove = e.AnimationQueue[0]
		e.AnimationQueue = e.AnimationQueue[1:]
		if e.CurrMove.SleepPre != 0 {
			time.Sleep(e.CurrMove.SleepPre)
		}
		e.tx = e.CurrMove.tx
		e.ty = e.CurrMove.ty
		vx := float64(e.tx - e.x)
		vy := float64(e.ty - e.y)
		speedVector := csg.NewVector(vx, vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
		e.vx = speedVector.X
		e.vy = speedVector.Y
	}
	e.x += e.vx
	e.y += e.vy
	if math.Abs(float64(e.tx-e.x))+math.Abs(float64(e.ty-e.y)) < 15 {
		if e.CurrMove != nil && e.CurrMove.DoneFunc != nil {
			if e.CurrMove.SleepPost != 0 {
				//time.Sleep(e.CurrMove.SleepPost)
			}
			e.CurrMove.DoneFunc()
		}
		if len(e.AnimationQueue) == 0 {
			e.x = e.tx
			e.y = e.ty
			e.vx = 0
			e.vy = 0
			e.CurrMove = nil
		} else {
			e.CurrMove = e.AnimationQueue[0]
			e.AnimationQueue = e.AnimationQueue[1:]
			if e.CurrMove.SleepPre != 0 {
				//time.Sleep(e.CurrMove.SleepPre)
			}
			e.tx = e.CurrMove.tx
			e.ty = e.CurrMove.ty
			vx := float64(e.tx - e.x)
			vy := float64(e.ty - e.y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
			e.vx = speedVector.X
			e.vy = speedVector.Y
		}

	}
}

func (m *MainGameState) Draw(screen *ebiten.Image) {

	// ebitenutil.DebugPrint(screen, "Hello, World!")
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(m.bgImage, op)
	op.GeoM.Reset()
	op.GeoM.Scale(0.6, 0.6)
	op.GeoM.Translate(ITEM_ICON_START_X, ITEM_ICON_START_Y)
	screen.DrawImage(m.ItemIcon, op)

	// m.defaultGamestate.MutexLock()
	//res := m.defaultGamestate.GetCurrentResource()
	//hp := m.defaultGamestate.GetCurrentHP()
	// m.defaultGamestate.MutexUnlock()
	m.mutex.Lock()
	text.Draw(screen, fmt.Sprintf("HP %d", m.hp), mplusDamage, 150, 40, color.RGBA{255, 0, 0, 255})
	if m.limiter != "" {
		ebitenutil.DebugPrint(screen, m.limiter)
	}
	text.Draw(screen, fmt.Sprintf("%d", m.combat), mplusResource, 380, 27, color.RGBA{255, 0, 0, 255})
	text.Draw(screen, fmt.Sprintf("%d", m.exploration), mplusResource, 380, 55, color.RGBA{0, 255, 0, 255})
	text.Draw(screen, fmt.Sprintf("%d", m.reputation), mplusResource, 500, 27, color.RGBA{127, 127, 0, 255})
	text.Draw(screen, fmt.Sprintf("%d", m.block), mplusResource, 500, 55, color.RGBA{127, 127, 0, 255})

	m.mutex.Unlock()

	op.GeoM.Reset()
	op.GeoM.Scale(0.35, 0.35)
	op.GeoM.Translate(350, 0)
	screen.DrawImage(m.iconCombat, op)
	// combat, ok := res.Detail[cards.RESOURCE_NAME_COMBAT]
	// if !ok {
	// 	combat = 0
	// }

	op.GeoM.Reset()
	op.GeoM.Scale(0.4, 0.4)
	op.GeoM.Translate(345, 30)
	screen.DrawImage(m.iconExplore, op)
	// explore, ok := res.Detail[cards.RESOURCE_NAME_EXPLORATION]
	// if !ok {
	// 	explore = 0
	// }

	op.GeoM.Reset()
	op.GeoM.Scale(0.15, 0.15)
	op.GeoM.Translate(450, 0)
	screen.DrawImage(m.Reputation, op)
	// rep, ok := res.Detail[cards.RESOURCE_NAME_REPUTATION]
	// if !ok {
	// 	rep = 0
	// }

	op.GeoM.Reset()
	op.GeoM.Scale(0.045, 0.045)
	op.GeoM.Translate(450, 30)
	screen.DrawImage(m.Block, op)
	// block, ok := res.Detail[cards.RESOURCE_NAME_BLOCK]
	// if !ok {
	// 	block = 0
	// }
	m.mutex.Lock()
	for _, c := range m.cardInHand {
		c.Draw(screen)
	}
	for _, c := range m.cardsPlayed {
		c.Draw(screen)
	}

	// fmt.Println(len(m.cardsInCenter))
	for _, c := range m.cardsInCenter {
		c.Draw(screen)
	}
	// center deck size
	m.mutex.Unlock()
	m.defaultGamestate.MutexLock()
	size := m.defaultGamestate.CardsInCenterDeck.Size()
	m.defaultGamestate.MutexUnlock()
	text.Draw(screen, fmt.Sprintf("%d", size), mplusResource, CENTER_DECK_START_X, CENTER_DECK_START_Y+50, color.RGBA{0, 255, 0, 255})
	op.GeoM.Reset()
	op.GeoM.Scale(HAND_SCALE, HAND_SCALE)
	op.GeoM.Translate(MAIN_DECK_X, MAIN_DECK_Y)
	screen.DrawImage(m.MainDeck, op)

	m.mutex.Lock()
	text.Draw(screen, fmt.Sprintf("%d", m.NumCardInDeck), mplusResource, MAIN_DECK_X+int(math.Floor(ORI_CARD_WIDTH*HAND_SCALE/2)), MAIN_DECK_Y+int(math.Floor(ORI_CARD_HEIGHT*HAND_SCALE*0.75)), color.RGBA{127, 127, 0, 255})
	m.mutex.Unlock()

	op.GeoM.Reset()
	op.GeoM.Scale(HAND_SCALE, HAND_SCALE)
	op.GeoM.Translate(DISCARD_START_X, DISCARD_START_Y)
	screen.DrawImage(m.DiscardPile, op)

	op.GeoM.Reset()
	// op.GeoM.Scale(HAND_SCALE, HAND_SCALE)
	op.GeoM.Translate(ENDTURN_START_X, ENDTURN_START_Y)
	screen.DrawImage(m.EndturnBtn, op)

	m.mutex.Lock()
	m.currentSubState.Draw(screen)
	for _, c := range m.cardsInLimbo {
		c.Draw(screen)
	}

	for _, c := range m.textInLimbo {
		c.Draw(screen)
	}
	m.mutex.Unlock()

	// if len(m.cardsPlayed) > 0 {
	// 	msg := fmt.Sprintf("Card1Pos=(%d,%d)\nCard1Target=(%d,%d)\nCard1V=(%d,%d)", m.cardsPlayed[0].x, m.cardsPlayed[0].y,
	// 		m.cardsPlayed[0].tx, m.cardsPlayed[0].ty, m.cardsPlayed[0].vx, m.cardsPlayed[0].vy)
	// 	ebitenutil.DebugPrint(screen, msg)
	// }

}
func (m *MainGameState) Update() error {
	dist := 0
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.dragMode {
		curX, _ := ebiten.CursorPosition()
		dist = curX - m.startDragX

	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// fmt.Println(inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft))
		m.dragMode = false
		if len(m.cardInHand) > 0 && m.cardInHand[0].x+float64(dist) > HAND_START_X {
			fmt.Println("terlalu ke kanan")
			dist = 0
			for indexCard, c := range m.cardInHand {
				c.x = math.Floor(HAND_START_X + float64(indexCard)*HAND_DIST_X)
			}
			// for idx, _ := range m.cardInHand {
			// 	m.cardInHand[idx].x_drag = 0
			// }
		}
		for _, c := range m.cardInHand {
			c.x += float64(dist)
			c.x_drag = 0
			// c.Update()
		}
	}
	// if m.dragMode {
	for _, c := range m.cardInHand {
		c.x_drag = dist
		c.Update()
	}
	// }

	for _, c := range m.cardsPlayed {
		c.Update()
	}

	for _, c := range m.cardsInCenter {

		c.Update()
	}

	newCardInLimbo := []*EbitenCard{}
	for _, c := range m.cardsInLimbo {
		c.Update()
		if c.tx != c.x || c.ty != c.y {
			newCardInLimbo = append(newCardInLimbo, c)
		}
	}
	m.cardsInLimbo = newCardInLimbo
	newTextInLimbo := []*EbitenText{}
	for _, c := range m.textInLimbo {
		c.Update()
		if c.tx != c.x || c.ty != c.y {
			newTextInLimbo = append(newTextInLimbo, c)
		}
	}
	// m.mutex.Lock()
	m.textInLimbo = newTextInLimbo
	// m.mutex.Unlock()

	return m.currentSubState.Update()
}
