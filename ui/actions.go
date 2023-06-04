package main

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"image/color"
	"math"
	"sync"
	"time"

	csg "github.com/kharism/golang-csg/core"
)

type OnDrawAction struct {
	mainGameState *MainGameState
}

func (d *OnDrawAction) DoAction(data map[string]interface{}) {
	// fmt.Println("OnDrawAction")
	drawnCards := data[cards.EVENT_ATTR_CARD_DRAWN].(cards.Card)

	newEbitenCard := NewEbitenCardFromCard(drawnCards)
	ll := mainGame.(*MainGameState)
	indexCard := len(ll.defaultGamestate.CardsInHand) - 1
	// fmt.Println("Draw card", drawnCards.GetName(), indexCard)
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
	fmt.Println("Draw card", drawnCards.GetName(), newEbitenCard.x, newEbitenCard.y, newEbitenCard.tx, newEbitenCard.ty)
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
		fmt.Println("MoveIds Play", moveIndex)
		newStartHand := HAND_START_X //float64(0.0)
		if len(newHand) > 0 {
			if moveIndex > 0 {
				newStartHand = newHand[0].x
			} else {

			}
		}
		fmt.Println("hand after played")
		for idx, c := range newHand {
			fmt.Println(idx, c.card.GetName(), c.x, c.tx)
		}
		for i := 0; i < len(newHand); i++ {
			fmt.Println("idx", i, newStartHand, HAND_DIST_X)
			newHand[i].tx = newStartHand + float64(i)*HAND_DIST_X //math.Floor(HAND_START_X + float64(i)*HAND_DIST_X)
			fmt.Println("idx", i, newHand[i].card.GetName(), newHand[i].tx, newHand[i].x)
			newHand[i].ty = HAND_START_Y
			vx := float64(newHand[i].tx - newHand[i].x)
			vy := float64(newHand[i].ty - newHand[i].y)
			if vx == 0 && vy == 0 {
				newHand[i].vx = 0
				newHand[i].vy = 0
			} else {
				speedVector := csg.NewVector(vx, vy, 0)
				speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)

				newHand[i].vx = speedVector.X
				newHand[i].vy = speedVector.Y
			}

		}
		mm.cardInHand = newHand
	}

	mm.mutex.Unlock()
	// fmt.Println(mm.cardInHand)
}

// handle EVENT_START_PUNISH
type onStartPunish struct {
	MainGameState *mainMainState
}

func (p *onStartPunish) DoAction(data map[string]interface{}) {

}

// handle EVENT_END_OF_TURN
type onEndOfTurnAction struct {
	MainGameState *mainMainState
}

func (p *onEndOfTurnAction) DoAction(data map[string]interface{}) {

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
			fmt.Println("Discard from hand", movedIdx)
			for i := movedIdx + 1; i < len(p.mainGameState.cardInHand); i++ {
				ebitenCard := sourceCard[i]
				ebitenCard.tx = ebitenCard.x - HAND_DIST_X
				vx := float64(ebitenCard.tx - ebitenCard.x)
				vy := float64(ebitenCard.ty - ebitenCard.y)
				speedVector := csg.NewVector(vx, vy, 0)
				speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
				ebitenCard.vx = speedVector.X
				ebitenCard.vy = speedVector.Y
				sourceCard[i] = ebitenCard
			}
		}
		fmt.Println("Done Discarding")
		for idx, c := range sourceCard {
			fmt.Println(idx, c.card.GetName(), c.x, c.tx)
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
	newAnim := &MoveAnimation{tx: ITEM_ICON_START_X, ty: ITEM_ICON_START_Y, Speed: CARD_MOVE_SPEED, SleepPre: 500 * time.Millisecond}
	ebitenCard.AnimationQueue = append(ebitenCard.AnimationQueue, newAnim)
	// ebitenCard.tx = ITEM_ICON_START_X
	// ebitenCard.ty = ITEM_ICON_START_Y
	// vx := float64(ebitenCard.tx - ebitenCard.x)
	// vy := float64(ebitenCard.ty - ebitenCard.y)
	// speedVector := csg.NewVector(vx, vy, 0)
	// speedVector = speedVector.Normalize().MultiplyScalar(CARD_MOVE_SPEED)
	// ebitenCard.vx = speedVector.X
	// ebitenCard.vy = speedVector.Y
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

type onTakeDamage struct {
	mainGameState *MainGameState
}

func (p *onTakeDamage) DoAction(data map[string]interface{}) {
	damageAmount := data[cards.EVENT_ATTR_CARD_TAKE_DAMAGE_AMMOUNT].(int)
	// TODO: add take damage/heal animation
	damageText := &EbitenText{text: fmt.Sprintf("%d", damageAmount), face: mplusResource, x: DMG_START_X, y: DMG_START_Y}
	if damageAmount > 0 {
		damageText.color = color.RGBA{255, 0, 0, 255}
	} else {
		damageText.color = color.RGBA{0, 255, 0, 255}
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)

	moveUp := &MoveAnimation{tx: damageText.x, ty: damageText.y - 50, Speed: 1}
	moveUp.DoneFunc = CreateDoneFunc(nil, wg)
	damageText.AnimationQueue = append(damageText.AnimationQueue, moveUp)
	p.mainGameState.textInLimbo = append(p.mainGameState.textInLimbo, damageText)
	wg.Wait()
}

type onPrePunish struct {
	mainGameState *MainGameState
}

func (p *onPrePunish) DoAction(data map[string]interface{}) {
	punishingCard := data[cards.EVENT_ATTR_BEFORE_PUNISH_CARD].(cards.Card)
	fmt.Println("Punishing cards", punishingCard.GetName())
	var animatedCard *EbitenCard
	for _, c := range p.mainGameState.cardsInCenter {
		if c.card == punishingCard {
			animatedCard = c
			break
		}
	}
	moveBack := &MoveAnimation{tx: animatedCard.x, ty: animatedCard.y - 20, Speed: 1}
	moveAtk := &MoveAnimation{tx: animatedCard.x, ty: animatedCard.y + 270, Speed: 10}
	moveReturn := &MoveAnimation{tx: animatedCard.x, ty: animatedCard.y, Speed: 5}
	// cc := make(chan string)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	moveReturn.DoneFunc = CreateDoneFunc(animatedCard, wg)
	moveQ := []*MoveAnimation{moveBack, moveAtk, moveReturn}
	animatedCard.AnimationQueue = append(animatedCard.AnimationQueue, moveQ...)
	// mutex2.Lock()
	wg.Wait()
}

type onCardStacked struct {
	mainGameState *MainGameState
}

func (p *onCardStacked) DoAction(data map[string]interface{}) {
	returnedCard := data[cards.EVENT_ATTR_CARD_STACKED].(cards.Card)
	fmt.Println("Stacking cards", returnedCard.GetName())
	source := data[cards.EVENT_ATTR_DISCARD_SOURCE].(string)
	if source == cards.DISCARD_SOURCE_HAND {
		moveIndex := -1
		newHandCard := []*EbitenCard{}
		for i := 0; i < len(p.mainGameState.cardInHand); i++ {
			if p.mainGameState.cardInHand[i].card == returnedCard {
				moveIndex = i
				ebitenCard := p.mainGameState.cardInHand[i]
				newAnim := &MoveAnimation{tx: MAIN_DECK_X, ty: MAIN_DECK_Y, Speed: CARD_MOVE_SPEED}
				ebitenCard.AnimationQueue = append(ebitenCard.AnimationQueue, newAnim)
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
			newAnim := &MoveAnimation{tx: math.Floor(HAND_START_X + float64(i)*HAND_DIST_X), ty: MAIN_DECK_Y, Speed: CARD_MOVE_SPEED}
			newHandCard[i].AnimationQueue = append(newHandCard[i].AnimationQueue, newAnim)
		}
		p.mainGameState.cardInHand = newHandCard

	} else if source == cards.DISCARD_SOURCE_NAN {
		ebitenCard := NewEbitenCardFromCard(returnedCard)
		ebitenCard.x = DISCARD_NA_SOURCE_X
		ebitenCard.y = DISCARD_NA_SOURCE_Y
		newAnim := &MoveAnimation{tx: MAIN_DECK_X, ty: MAIN_DECK_Y, Speed: CARD_MOVE_SPEED, SleepPre: 750 * time.Millisecond}
		ebitenCard.AnimationQueue = append(ebitenCard.AnimationQueue, newAnim)
		p.mainGameState.cardsInLimbo = append(p.mainGameState.cardsInLimbo, ebitenCard)
	}
}
