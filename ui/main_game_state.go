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

	csg "github.com/kharism/golang-csg/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
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
	// cards in limbo meaning cards that is moving into cooldownpile or banished pile
	// they have still visible until they reach those position
	cardsInLimbo     []*EbitenCard
	stateChanger     AbstractStateChanger
	detailViewCard   *EbitenCard
	mutex            *sync.Mutex
	defaultGamestate *gamestate.DefaultGamestate
	limiter          string

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
	m *MainGameState
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
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		xCurInt, yCurInt := ebiten.CursorPosition()
		xCur, yCur := float64(xCurInt), float64(yCurInt)

		if xCur > DISCARD_START_X && xCur < DISCARD_START_X+HAND_SCALE*ORI_CARD_WIDTH {
			s.m.cardListState.cards = s.m.defaultGamestate.CardsDiscarded.List()
			s.m.currentSubState = s.m.cardListState
		} else if xCur > ENDTURN_START_X {
			fmt.Println("Endturn")
			go func() {
				s.m.defaultGamestate.EndTurn()
				s.m.defaultGamestate.BeginTurn()
			}()
		} else if yCur > HAND_START_Y && xCur < DISCARD_START_X && xCur >= HAND_START_X {
			// left click on hand
			for i := len(s.m.cardInHand) - 1; i >= 0; i-- {
				if s.m.cardInHand[i].x < xCur {
					go s.m.defaultGamestate.PlayCard(s.m.cardInHand[i].card)
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
							go s.m.defaultGamestate.DefeatCard(clickedCard.card)
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
	s.prevSubState = s.m.currentSubState
	s.m.currentSubState = s
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

type OnDrawAction struct {
	mainGameState *MainGameState
}

const (
	CARD_MOVE_SPEED = 10
)

func (d *OnDrawAction) DoAction(data map[string]interface{}) {
	// fmt.Println("OnDrawAction")
	drawnCards := data[cards.EVENT_ATTR_CARD_DRAWN].(cards.Card)

	newEbitenCard := NewEbitenCardFromCard(drawnCards)
	ll := mainGame.(*MainGameState)
	indexCard := len(ll.defaultGamestate.CardsInHand) - 1
	fmt.Println("Draw card", drawnCards.GetName(), indexCard)
	newEbitenCard.x = math.Floor(MAIN_DECK_X)
	newEbitenCard.y = math.Floor(MAIN_DECK_Y)
	newEbitenCard.tx = math.Floor(HAND_START_X + float64(indexCard)*HAND_DIST_X)
	newEbitenCard.ty = HAND_START_Y
	vx := float64(newEbitenCard.tx - newEbitenCard.x)
	vy := float64(newEbitenCard.ty - newEbitenCard.y)
	speedVector := csg.NewVector(vx, vy, 0)
	speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
	newEbitenCard.vx = speedVector.X
	newEbitenCard.vy = speedVector.Y
	// fmt.Println(newEbitenCard.x, newEbitenCard.y, newEbitenCard.tx, newEbitenCard.ty)
	ll.mutex.Lock()
	ll.cardInHand = append(ll.cardInHand, newEbitenCard)
	ll.mutex.Unlock()
}

type OnPlayAction struct {
	mainGameState *MainGameState
}

func (p *OnPlayAction) DoAction(data map[string]interface{}) {
	playedCards := data[cards.EVENT_ATTR_CARD_PLAYED].(cards.Card)
	fmt.Println("Play Card", playedCards.GetName())
	mm := mainGame.(*MainGameState)
	newHand := []*EbitenCard{}
	mm.mutex.Lock()
	tx := PLAYED_START_X + 45*len(mm.cardsPlayed)
	ty := PLAYED_START_Y
	moveIndex := -1
	for idx, val := range mm.cardInHand {
		txOld, tyOld := val.tx, val.ty
		if val.card == playedCards {
			moveIndex = idx
			val.tx = float64(tx)
			val.ty = float64(ty)
			vx := float64(val.tx - val.x)
			vy := float64(val.ty - val.y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			val.vx = speedVector.X
			val.vy = speedVector.Y

			mm.cardsPlayed = append(mm.cardsPlayed, val)
		} else {
			newHand = append(newHand, val)
		}
		fmt.Printf("%s old (%f,%f) Target (%f,%f)\n", val.card.GetName(), txOld, tyOld, val.tx, val.ty)

	}
	if moveIndex == -1 {
		fmt.Println(playedCards, playedCards.GetName())
	}
	// move any card on the right of our picked cards
	if moveIndex > -1 {
		for i := moveIndex; i < len(newHand); i++ {
			newHand[i].tx = math.Floor(HAND_START_X + float64(i)*HAND_DIST_X)
			newHand[i].ty = HAND_START_Y
			vx := float64(newHand[i].tx - newHand[i].x)
			vy := float64(newHand[i].ty - newHand[i].y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			newHand[i].vx = speedVector.X
			newHand[i].vy = speedVector.Y
		}
		mm.cardInHand = newHand
	}

	mm.mutex.Unlock()
	// fmt.Println(mm.cardInHand)
}

type onBanishAction struct {
	mainGameState *MainGameState
}

func (p *onBanishAction) DoAction(data map[string]interface{}) {
	cardDiscarded := data[cards.EVENT_ATTR_CARD_BANISHED].(cards.Card)
	source := data[cards.EVENT_ATTR_DISCARD_SOURCE].(string)
	fmt.Println("Discard card", cardDiscarded.GetName())
	defer p.mainGameState.mutex.Unlock()
	p.mainGameState.mutex.Lock()
	sourceCard := []*EbitenCard{}
	newSource := []*EbitenCard{}
	if source == cards.DISCARD_SOURCE_HAND {
		sourceCard = p.mainGameState.cardInHand
	} else if source == cards.DISCARD_SOURCE_PLAYED {
		// newPlayed := []*EbitenCard{}
		sourceCard = p.mainGameState.cardsPlayed
		// p.mainGameState.cardsPlayed = newPlayed
	} else if source == cards.DISCARD_SOURCE_CENTER {
		sourceCard = p.mainGameState.cardsInCenter
	}

	movedIdx := -1
	for i := 0; i < len(sourceCard); i++ {
		if sourceCard[i].card == cardDiscarded {
			movedIdx = i
			ebitenCard := sourceCard[i]
			ebitenCard.tx = BANISHED_START_X
			ebitenCard.ty = BANISHED_START_Y
			vx := float64(ebitenCard.tx - ebitenCard.x)
			vy := float64(ebitenCard.ty - ebitenCard.y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			ebitenCard.vx = speedVector.X
			ebitenCard.vy = speedVector.Y
			// fmt.Println("DetectDiscarded", cardDiscarded.GetName(), source, ebitenCard.vx, ebitenCard.vy)
			p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
		} else {
			newSource = append(newSource, sourceCard[i])
		}
	}
	if source == cards.DISCARD_SOURCE_HAND {
		p.mainGameState.cardInHand = newSource
		// move cards on the right side to left
		if len(p.mainGameState.cardInHand) > 0 {
			for i := movedIdx; i < len(p.mainGameState.cardInHand); i++ {
				ebitenCard := sourceCard[i]
				ebitenCard.tx -= HAND_DIST_X
				vx := float64(ebitenCard.tx - ebitenCard.x)
				vy := float64(ebitenCard.ty - ebitenCard.y)
				speedVector := csg.NewVector(vx, vy, 0)
				speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
				ebitenCard.vx = speedVector.X
				ebitenCard.vy = speedVector.Y
				sourceCard[i] = ebitenCard
			}
		}
		// sourceCard = p.mainGameState.cardInHand
	} else if source == cards.DISCARD_SOURCE_CENTER {
		p.mainGameState.cardsInCenter = newSource
		if len(p.mainGameState.cardsInCenter) > 0 && movedIdx != -1 {
			for i := movedIdx; i < len(p.mainGameState.cardsInCenter); i++ {
				ebitenCard := sourceCard[i]
				ebitenCard.tx -= HAND_DIST_X
				vx := float64(ebitenCard.tx - ebitenCard.x)
				vy := float64(ebitenCard.ty - ebitenCard.y)
				speedVector := csg.NewVector(vx, vy, 0)
				speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
				ebitenCard.vx = speedVector.X
				ebitenCard.vy = speedVector.Y
				sourceCard[i] = ebitenCard
			}
		}
	} else if source == cards.DISCARD_SOURCE_PLAYED {
		// newPlayed := []*EbitenCard{}
		// sourceCard = p.mainGameState.cardsPlayed
		p.mainGameState.cardsPlayed = newSource
	}
}

type onDiscardAction struct {
	mainGameState *MainGameState
}

func (p *onDiscardAction) DoAction(data map[string]interface{}) {
	cardDiscarded := data[cards.EVENT_ATTR_CARD_DISCARDED].(cards.Card)
	source := data[cards.EVENT_ATTR_DISCARD_SOURCE].(string)
	fmt.Println("Discard card", cardDiscarded.GetName())
	defer p.mainGameState.mutex.Unlock()
	p.mainGameState.mutex.Lock()
	// var ebitenCard *EbitenCard
	sourceCard := []*EbitenCard{}
	newSource := []*EbitenCard{}
	if source == cards.DISCARD_SOURCE_HAND {
		// p.mainGameState.cardInHand = newHand
		sourceCard = p.mainGameState.cardInHand
	} else if source == cards.DISCARD_SOURCE_PLAYED {
		// newPlayed := []*EbitenCard{}
		sourceCard = p.mainGameState.cardsPlayed
		// p.mainGameState.cardsPlayed = newPlayed
	} else if source == cards.DISCARD_SOURCE_NAN {
		ebitenCard := NewEbitenCardFromCard(cardDiscarded)
		ebitenCard.tx = DISCARD_START_X
		ebitenCard.ty = DISCARD_START_Y
		ebitenCard.x = DISCARD_NA_SOURCE_X
		ebitenCard.y = DISCARD_NA_SOURCE_Y
		vx := float64(ebitenCard.tx - ebitenCard.x)
		vy := float64(ebitenCard.ty - ebitenCard.y)
		speedVector := csg.NewVector(vx, vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
		ebitenCard.vx = speedVector.X
		ebitenCard.vy = speedVector.Y
		p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
	}
	movedIdx := -1
	fmt.Println("SourceCard", len(sourceCard))
	for i := 0; i < len(sourceCard); i++ {
		if sourceCard[i].card == cardDiscarded {
			movedIdx = i
			ebitenCard := sourceCard[i]
			ebitenCard.tx = DISCARD_START_X
			ebitenCard.ty = DISCARD_START_Y
			vx := float64(ebitenCard.tx - ebitenCard.x)
			vy := float64(ebitenCard.ty - ebitenCard.y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			ebitenCard.vx = speedVector.X
			ebitenCard.vy = speedVector.Y
			// fmt.Println("DetectDiscarded", cardDiscarded.GetName(), source, ebitenCard.vx, ebitenCard.vy)
			p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
		} else {
			newSource = append(newSource, sourceCard[i])
		}
	}
	if source == cards.DISCARD_SOURCE_HAND {
		p.mainGameState.cardInHand = newSource
		// move cards on the right side to left
		if movedIdx == -1 {
			// debug this crap
			fmt.Println("Card Not found", cardDiscarded.GetName())

		}
		if len(p.mainGameState.cardInHand) > 0 {
			for i := movedIdx; i < len(p.mainGameState.cardInHand); i++ {
				ebitenCard := sourceCard[i]
				ebitenCard.tx -= HAND_DIST_X
				vx := float64(ebitenCard.tx - ebitenCard.x)
				vy := float64(ebitenCard.ty - ebitenCard.y)
				speedVector := csg.NewVector(vx, vy, 0)
				speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
				ebitenCard.vx = speedVector.X
				ebitenCard.vy = speedVector.Y
				sourceCard[i] = ebitenCard
			}
		}
		// sourceCard = p.mainGameState.cardInHand
	} else if source == cards.DISCARD_SOURCE_PLAYED {
		// newPlayed := []*EbitenCard{}
		// sourceCard = p.mainGameState.cardsPlayed
		p.mainGameState.cardsPlayed = newSource
	}

}

type onCenterDrawAction struct {
	mainGameState *MainGameState
}

func (p *onCenterDrawAction) DoAction(data map[string]interface{}) {

	drawnCards := data[cards.EVENT_ATTR_CARD_DRAWN].(cards.Card)
	newEbitenCard := NewEbitenCardFromCard(drawnCards)
	fmt.Println("center Draw", drawnCards.GetName())
	isDisarmedTrap := false
	isATrap := false
	if _, ok := drawnCards.(cards.Trapper); ok {
		j := drawnCards.(cards.Trapper)
		isATrap = true
		if j.IsDisarmed() {
			isDisarmedTrap = true
		}
	}
	ll := mainGame.(*MainGameState)
	indexCard := len(ll.defaultGamestate.CenterCards)
	newEbitenCard.x = math.Floor(CENTER_DECK_START_X)
	newEbitenCard.y = math.Floor(CENTER_DECK_START_Y)
	if isDisarmedTrap {
		newEbitenCard.tx = BANISHED_START_X
		newEbitenCard.ty = BANISHED_START_Y
	} else {
		newEbitenCard.tx = math.Floor(CENTER_START_X + float64(indexCard)*HAND_DIST_X)
		newEbitenCard.ty = CENTER_START_Y
	}

	vx := float64(newEbitenCard.tx - newEbitenCard.x)
	vy := float64(newEbitenCard.ty - newEbitenCard.y)
	speedVector := csg.NewVector(vx, vy, 0)
	speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
	newEbitenCard.vx = speedVector.X
	newEbitenCard.vy = speedVector.Y
	if isATrap {
		fmt.Println(newEbitenCard.x, newEbitenCard.y, newEbitenCard.tx, newEbitenCard.ty, newEbitenCard.vx, newEbitenCard.vy)
	}

	ll.mutex.Lock()
	if isDisarmedTrap {
		fmt.Println("Append disarmed trap in limbo")
		ll.cardsInLimbo = append(ll.cardsInLimbo, newEbitenCard)
	} else {
		fmt.Println("append", drawnCards.GetName())
		ll.cardsInCenter = append(ll.cardsInCenter, newEbitenCard)
		fmt.Println("====", len(ll.cardsInCenter))
	}

	ll.mutex.Unlock()

	// if we draw the last card on center deck, trigger you win
	// TODO: add alternative win condition
	// if ll.defaultGamestate.CardsInCenterDeck.Size() == 0 {
	// 	ll.currentSubState = ll.gameoverState
	// }
}

type onExplorationAction struct {
	mainGameState *MainGameState
}

func (p *onExplorationAction) DoAction(data map[string]interface{}) {
	exploredCard := data[cards.EVENT_ATTR_CARD_EXPLORED].(cards.Card)
	newCenterCard := []*EbitenCard{}
	fmt.Println("Explore", exploredCard.GetName())
	defer p.mainGameState.mutex.Unlock()
	p.mainGameState.mutex.Lock()
	moveIndex := -1
	for i := 0; i < len(p.mainGameState.cardsInCenter); i++ {
		if p.mainGameState.cardsInCenter[i].card == exploredCard {
			moveIndex = i
			ebitenCard := p.mainGameState.cardsInCenter[i]
			ebitenCard.tx = BANISHED_START_X
			ebitenCard.ty = BANISHED_START_Y
			vx := float64(ebitenCard.tx - ebitenCard.x)
			vy := float64(ebitenCard.ty - ebitenCard.y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			ebitenCard.vx = speedVector.X
			ebitenCard.vy = speedVector.Y
			p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
		} else {
			newCenterCard = append(newCenterCard, p.mainGameState.cardsInCenter[i])
		}
	}
	if moveIndex != -1 {
		for i := moveIndex; i < len(newCenterCard); i++ {
			newCenterCard[i].tx = math.Floor(CENTER_START_X + float64(i)*HAND_DIST_X)
			newCenterCard[i].ty = CENTER_START_Y
			vx := float64(newCenterCard[i].tx - newCenterCard[i].x)
			vy := float64(newCenterCard[i].ty - newCenterCard[i].y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			newCenterCard[i].vx = speedVector.X
			newCenterCard[i].vy = speedVector.Y
		}
	}

	p.mainGameState.cardsInCenter = newCenterCard
}

type onDefeatAction struct {
	mainGameState *MainGameState
}

func (p *onDefeatAction) DoAction(data map[string]interface{}) {
	fmt.Println("defeat")
	exploredCard := data[cards.EVENT_ATTR_CARD_DEFEATED].(cards.Card)
	newCenterCard := []*EbitenCard{}
	defer p.mainGameState.mutex.Unlock()
	p.mainGameState.mutex.Lock()
	moveIndex := -1
	for i := 0; i < len(p.mainGameState.cardsInCenter); i++ {
		if p.mainGameState.cardsInCenter[i].card == exploredCard {
			moveIndex = i
			ebitenCard := p.mainGameState.cardsInCenter[i]
			ebitenCard.tx = BANISHED_START_X
			ebitenCard.ty = BANISHED_START_Y
			vx := float64(ebitenCard.tx - ebitenCard.x)
			vy := float64(ebitenCard.ty - ebitenCard.y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			ebitenCard.vx = speedVector.X
			ebitenCard.vy = speedVector.Y
			p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
		} else {
			newCenterCard = append(newCenterCard, p.mainGameState.cardsInCenter[i])
		}
	}
	if moveIndex != -1 {
		for i := moveIndex; i < len(newCenterCard); i++ {
			newCenterCard[i].tx = math.Floor(CENTER_START_X + float64(i)*HAND_DIST_X)
			newCenterCard[i].ty = CENTER_START_Y
			vx := float64(newCenterCard[i].tx - newCenterCard[i].x)
			vy := float64(newCenterCard[i].ty - newCenterCard[i].y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			newCenterCard[i].vx = speedVector.X
			newCenterCard[i].vy = speedVector.Y
		}
	}

	p.mainGameState.cardsInCenter = newCenterCard
}

type onItemAdd struct {
	mainGameState *MainGameState
}

func (p *onItemAdd) DoAction(data map[string]interface{}) {
	addedItem := data[cards.EVENT_ATTR_ITEM_ADDED].(cards.Card)
	ebitenCard := NewEbitenCardFromCard(addedItem)
	ebitenCard.x = DISCARD_NA_SOURCE_X
	ebitenCard.y = DISCARD_NA_SOURCE_Y
	ebitenCard.tx = ITEM_ICON_START_X
	ebitenCard.ty = ITEM_ICON_START_Y
	vx := float64(ebitenCard.tx - ebitenCard.x)
	vy := float64(ebitenCard.ty - ebitenCard.y)
	speedVector := csg.NewVector(vx, vy, 0)
	speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
	ebitenCard.vx = speedVector.X
	ebitenCard.vy = speedVector.Y
	p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
}

type onDisarmAction struct {
	mainGameState *MainGameState
}

func (p *onDisarmAction) DoAction(data map[string]interface{}) {
	// fmt.Println("defeat")
	exploredCard := data[cards.EVENT_ATTR_TRAP_REMOVED].(cards.Card)
	newCenterCard := []*EbitenCard{}
	defer p.mainGameState.mutex.Unlock()
	p.mainGameState.mutex.Lock()
	moveIndex := -1
	for i := 0; i < len(p.mainGameState.cardsInCenter); i++ {
		if p.mainGameState.cardsInCenter[i].card == exploredCard {
			moveIndex = i
			ebitenCard := p.mainGameState.cardsInCenter[i]
			ebitenCard.tx = BANISHED_START_X
			ebitenCard.ty = BANISHED_START_Y
			vx := float64(ebitenCard.tx - ebitenCard.x)
			vy := float64(ebitenCard.ty - ebitenCard.y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			ebitenCard.vx = speedVector.X
			ebitenCard.vy = speedVector.Y
			p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
		} else {
			newCenterCard = append(newCenterCard, p.mainGameState.cardsInCenter[i])
		}
	}
	if moveIndex != -1 {
		for i := moveIndex; i < len(newCenterCard); i++ {
			newCenterCard[i].tx = math.Floor(CENTER_START_X + float64(i)*HAND_DIST_X)
			newCenterCard[i].ty = CENTER_START_Y
			vx := float64(newCenterCard[i].tx - newCenterCard[i].x)
			vy := float64(newCenterCard[i].ty - newCenterCard[i].y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			newCenterCard[i].vx = speedVector.X
			newCenterCard[i].vy = speedVector.Y
		}
	}

	p.mainGameState.cardsInCenter = newCenterCard
}

type onRecruitAction struct {
	mainGameState *MainGameState
}

func (p *onRecruitAction) DoAction(data map[string]interface{}) {
	fmt.Println("Recruit")
	exploredCard := data[cards.EVENT_ATTR_CARD_RECRUITED].(cards.Card)
	newCenterCard := []*EbitenCard{}
	defer p.mainGameState.mutex.Unlock()
	p.mainGameState.mutex.Lock()
	moveIndex := -1
	for i := 0; i < len(p.mainGameState.cardsInCenter); i++ {
		if p.mainGameState.cardsInCenter[i].card == exploredCard {
			moveIndex = i
			ebitenCard := p.mainGameState.cardsInCenter[i]
			ebitenCard.tx = DISCARD_START_X
			ebitenCard.ty = DISCARD_START_Y
			vx := float64(ebitenCard.tx - ebitenCard.x)
			vy := float64(ebitenCard.ty - ebitenCard.y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			ebitenCard.vx = speedVector.X
			ebitenCard.vy = speedVector.Y
			p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
		} else {
			newCenterCard = append(newCenterCard, p.mainGameState.cardsInCenter[i])
		}
	}
	fmt.Println("Geser Recruit")
	for i := moveIndex; i < len(newCenterCard); i++ {
		newCenterCard[i].tx = math.Floor(CENTER_START_X + float64(i)*HAND_DIST_X)
		newCenterCard[i].ty = CENTER_START_Y
		vx := float64(newCenterCard[i].tx - newCenterCard[i].x)
		vy := float64(newCenterCard[i].ty - newCenterCard[i].y)
		speedVector := csg.NewVector(vx, vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
		newCenterCard[i].vx = speedVector.X
		newCenterCard[i].vy = speedVector.Y
		// fmt.Sprintf("%d %f %f\n", i, newCenterCard[i].tx, newCenterCard[i].ty)
	}
	p.mainGameState.cardsInCenter = newCenterCard
}

type onGotoCenterDeckAction struct {
	mainGameState *MainGameState
}

func (p *onGotoCenterDeckAction) DoAction(data map[string]interface{}) {
	returnedCard := data[cards.EVENT_ATTR_CARD_GOTO_CENTER].(cards.Card)
	source := data[cards.EVENT_ATTR_DISCARD_SOURCE].(string)
	fmt.Println("Center stuff", returnedCard.GetName(), source)
	defer p.mainGameState.mutex.Unlock()
	p.mainGameState.mutex.Lock()
	if source == cards.DISCARD_SOURCE_CENTER {
		moveIndex := -1
		newCenterCard := []*EbitenCard{}
		for i := 0; i < len(p.mainGameState.cardsInCenter); i++ {
			if p.mainGameState.cardsInCenter[i].card == returnedCard {
				moveIndex = i
				ebitenCard := p.mainGameState.cardsInCenter[i]
				ebitenCard.tx = CENTER_DECK_START_X
				ebitenCard.ty = CENTER_DECK_START_Y
				vx := float64(ebitenCard.tx - ebitenCard.x)
				vy := float64(ebitenCard.ty - ebitenCard.y)
				speedVector := csg.NewVector(vx, vy, 0)
				speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
				ebitenCard.vx = speedVector.X
				ebitenCard.vy = speedVector.Y
				p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
			} else {
				newCenterCard = append(newCenterCard, p.mainGameState.cardsInCenter[i])
			}
		}
		for i := moveIndex; i < len(newCenterCard); i++ {
			newCenterCard[i].tx = math.Floor(CENTER_START_X + float64(i)*HAND_DIST_X)
			newCenterCard[i].ty = CENTER_START_Y
			vx := float64(newCenterCard[i].tx - newCenterCard[i].x)
			vy := float64(newCenterCard[i].ty - newCenterCard[i].y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			newCenterCard[i].vx = speedVector.X
			newCenterCard[i].vy = speedVector.Y
			// fmt.Sprintf("%d %f %f\n", i, newCenterCard[i].tx, newCenterCard[i].ty)
		}
		p.mainGameState.cardsInCenter = newCenterCard

	} else if source == cards.DISCARD_SOURCE_NAN {
		ebitenCard := NewEbitenCardFromCard(returnedCard)
		ebitenCard.x = DISCARD_NA_SOURCE_X
		ebitenCard.y = DISCARD_NA_SOURCE_Y
		ebitenCard.tx = CENTER_DECK_START_X
		ebitenCard.ty = CENTER_DECK_START_Y
		vx := float64(ebitenCard.tx - ebitenCard.x)
		vy := float64(ebitenCard.ty - ebitenCard.y)
		speedVector := csg.NewVector(vx, vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
		ebitenCard.vx = speedVector.X
		ebitenCard.vy = speedVector.Y
		p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
	}
}

type onCardStacked struct {
	mainGameState *MainGameState
}

func (p *onCardStacked) DoAction(data map[string]interface{}) {
	returnedCard := data[cards.EVENT_ATTR_CARD_STACKED].(cards.Card)
	source := data[cards.EVENT_ATTR_DISCARD_SOURCE].(string)
	if source == cards.DISCARD_SOURCE_HAND {
		moveIndex := -1
		newHandCard := []*EbitenCard{}
		for i := 0; i < len(p.mainGameState.cardInHand); i++ {
			if p.mainGameState.cardInHand[i].card == returnedCard {
				moveIndex = i
				ebitenCard := p.mainGameState.cardInHand[i]
				ebitenCard.tx = MAIN_DECK_X
				ebitenCard.ty = MAIN_DECK_Y
				vx := float64(ebitenCard.tx - ebitenCard.x)
				vy := float64(ebitenCard.ty - ebitenCard.y)
				speedVector := csg.NewVector(vx, vy, 0)
				speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
				ebitenCard.vx = speedVector.X
				ebitenCard.vy = speedVector.Y
				p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
			} else {
				newHandCard = append(newHandCard, p.mainGameState.cardInHand[i])
			}
		}
		fmt.Println("Move index", moveIndex)
		if moveIndex == -1 {
			return
		}
		for i := moveIndex; i < len(newHandCard); i++ {
			newHandCard[i].tx = math.Floor(HAND_START_X + float64(i)*HAND_DIST_X)
			newHandCard[i].ty = MAIN_DECK_Y
			vx := float64(newHandCard[i].tx - newHandCard[i].x)
			vy := float64(newHandCard[i].ty - newHandCard[i].y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
			newHandCard[i].vx = speedVector.X
			newHandCard[i].vy = speedVector.Y
			// fmt.Sprintf("%d %f %f\n", i, newCenterCard[i].tx, newCenterCard[i].ty)
		}
		p.mainGameState.cardInHand = newHandCard

	} else if source == cards.DISCARD_SOURCE_NAN {
		ebitenCard := NewEbitenCardFromCard(returnedCard)
		ebitenCard.x = DISCARD_NA_SOURCE_X
		ebitenCard.y = DISCARD_NA_SOURCE_Y
		ebitenCard.tx = MAIN_DECK_X
		ebitenCard.ty = MAIN_DECK_Y
		vx := float64(ebitenCard.tx - ebitenCard.x)
		vy := float64(ebitenCard.ty - ebitenCard.y)
		speedVector := csg.NewVector(vx, vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
		ebitenCard.vx = speedVector.X
		ebitenCard.vy = speedVector.Y
		p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
	}
}

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
	mainState := &mainMainState{m: mgs}
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

func (m *MainGameState) Draw(screen *ebiten.Image) {

	// ebitenutil.DebugPrint(screen, "Hello, World!")
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(m.bgImage, op)
	op.GeoM.Reset()
	op.GeoM.Scale(0.6, 0.6)
	op.GeoM.Translate(ITEM_ICON_START_X, ITEM_ICON_START_Y)
	screen.DrawImage(m.ItemIcon, op)

	res := m.defaultGamestate.GetCurrentResource()
	hp := m.defaultGamestate.GetCurrentHP()
	text.Draw(screen, fmt.Sprintf("HP %d", hp), mplusResource, 150, 40, color.RGBA{255, 0, 0, 255})
	if m.limiter != "" {
		ebitenutil.DebugPrint(screen, m.limiter)
	}

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

	op.GeoM.Reset()
	op.GeoM.Scale(0.3, 0.3)
	op.GeoM.Translate(750, -10)
	screen.DrawImage(m.Reputation, op)
	rep, ok := res.Detail[cards.RESOURCE_NAME_REPUTATION]
	if !ok {
		rep = 0
	}
	text.Draw(screen, fmt.Sprintf("%d", rep), mplusResource, 850, 40, color.RGBA{127, 127, 0, 255})
	op.GeoM.Reset()
	op.GeoM.Scale(0.09, 0.09)
	op.GeoM.Translate(900, -10)
	screen.DrawImage(m.Block, op)
	block, ok := res.Detail[cards.RESOURCE_NAME_BLOCK]
	if !ok {
		block = 0
	}
	text.Draw(screen, fmt.Sprintf("%d", block), mplusResource, 960, 40, color.RGBA{127, 127, 0, 255})

	for _, c := range m.cardInHand {
		c.Draw(screen)
	}
	for _, c := range m.cardsPlayed {
		c.Draw(screen)
	}
	for _, c := range m.cardsInLimbo {
		c.Draw(screen)
	}
	// fmt.Println(len(m.cardsInCenter))
	for _, c := range m.cardsInCenter {
		c.Draw(screen)
	}
	// center deck size
	size := m.defaultGamestate.CardsInCenterDeck.Size()
	text.Draw(screen, fmt.Sprintf("%d", size), mplusResource, CENTER_DECK_START_X, CENTER_DECK_START_Y+50, color.RGBA{0, 255, 0, 255})
	op.GeoM.Reset()
	op.GeoM.Scale(HAND_SCALE, HAND_SCALE)
	op.GeoM.Translate(MAIN_DECK_X, MAIN_DECK_Y)
	screen.DrawImage(m.MainDeck, op)

	op.GeoM.Reset()
	op.GeoM.Scale(HAND_SCALE, HAND_SCALE)
	op.GeoM.Translate(DISCARD_START_X, DISCARD_START_Y)
	screen.DrawImage(m.DiscardPile, op)

	op.GeoM.Reset()
	// op.GeoM.Scale(HAND_SCALE, HAND_SCALE)
	op.GeoM.Translate(ENDTURN_START_X, ENDTURN_START_Y)
	screen.DrawImage(m.EndturnBtn, op)

	m.currentSubState.Draw(screen)

	// if len(m.cardsPlayed) > 0 {
	// 	msg := fmt.Sprintf("Card1Pos=(%d,%d)\nCard1Target=(%d,%d)\nCard1V=(%d,%d)", m.cardsPlayed[0].x, m.cardsPlayed[0].y,
	// 		m.cardsPlayed[0].tx, m.cardsPlayed[0].ty, m.cardsPlayed[0].vx, m.cardsPlayed[0].vy)
	// 	ebitenutil.DebugPrint(screen, msg)
	// }

}
func (m *MainGameState) Update() error {

	return m.currentSubState.Update()
}
