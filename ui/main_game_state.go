package main

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/gamestate"
	"image/color"
	"log"
	"math"
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

	// sub-states
	currentSubState SubState
	mainState       *mainMainState
	detailState     *detailState
	cardPicker      *cardPickState
	cardListState   *cardListState
	gameoverState   *gameOverSubstate
}
type SubState interface {
	Draw(screen *ebiten.Image)
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
		s.m.detailState.prevSubState = s
		s.m.currentSubState = s.m.detailState
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
		} else if yCur > HAND_START_Y && xCur < DISCARD_START_X {
			// left click on hand
			for i := len(s.m.cardInHand) - 1; i >= 0; i-- {
				if s.m.cardInHand[i].x < xCur {
					go s.m.defaultGamestate.PlayCard(s.m.cardInHand[i].card)
					//s.m.detailViewCard = s.m.cardInHand[i]
					//fmt.Println("cardIndex at", i)
					break
				}
			}
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
						go s.m.defaultGamestate.DefeatCard(clickedCard.card)
					}
					//go s.m.defaultGamestate.PlayCard(s.m.cardInHand[i].card)
					//s.m.detailViewCard = s.m.cardInHand[i]
					//fmt.Println("cardIndex at", i)
					break
				}
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
		// s.m.currentSubState = s.m.cardPicker
		s.m.defaultGamestate.Draw()
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
	// move any card on the right of our picked cards
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
	mm.mutex.Unlock()
	// fmt.Println(mm.cardInHand)
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
	}
	movedIdx := -1
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
	// fmt.Println("center Draw", drawnCards.GetName())
	ll := mainGame.(*MainGameState)
	indexCard := len(ll.defaultGamestate.CenterCards)
	newEbitenCard.x = math.Floor(CENTER_DECK_START_X)
	newEbitenCard.y = math.Floor(CENTER_DECK_START_Y)
	newEbitenCard.tx = math.Floor(CENTER_START_X + float64(indexCard)*HAND_DIST_X)
	newEbitenCard.ty = CENTER_START_Y
	vx := float64(newEbitenCard.tx - newEbitenCard.x)
	vy := float64(newEbitenCard.ty - newEbitenCard.y)
	speedVector := csg.NewVector(vx, vy, 0)
	speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
	newEbitenCard.vx = speedVector.X
	newEbitenCard.vy = speedVector.Y
	// fmt.Println(newEbitenCard.x, newEbitenCard.y, newEbitenCard.tx, newEbitenCard.ty, newEbitenCard.vx, newEbitenCard.vy)
	ll.mutex.Lock()
	ll.cardsInCenter = append(ll.cardsInCenter, newEbitenCard)
	ll.mutex.Unlock()

	// if we draw the last card on center deck, trigger you win
	// TODO: add alternative win condition
	if ll.defaultGamestate.CardsInCenterDeck.Size() == 0 {
		ll.currentSubState = ll.gameoverState
	}
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
	}
	mainState := &mainMainState{m: mgs}
	detailState := &detailState{m: mgs}
	cardpicker := &cardPickState{m: mgs, pickedCards: make(chan int)}
	cardListState := &cardListState{m: mgs}
	mgs.gameoverState = &gameOverSubstate{m: mgs}
	mgs.currentSubState = mainState
	mgs.mainState = mainState
	mgs.detailState = detailState
	mgs.cardPicker = cardpicker
	mgs.cardListState = cardListState
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

	if len(m.cardsPlayed) > 0 {
		msg := fmt.Sprintf("Card1Pos=(%d,%d)\nCard1Target=(%d,%d)\nCard1V=(%d,%d)", m.cardsPlayed[0].x, m.cardsPlayed[0].y,
			m.cardsPlayed[0].tx, m.cardsPlayed[0].ty, m.cardsPlayed[0].vx, m.cardsPlayed[0].vy)
		ebitenutil.DebugPrint(screen, msg)
	}

}
func (m *MainGameState) Update() error {
	for _, c := range m.cardInHand {
		c.Update()
	}
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

	return nil
}
