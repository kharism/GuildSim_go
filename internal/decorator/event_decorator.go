package decorator

import (
	"github/kharism/GuildSim_go/internal/cards"
	mrand "math/rand"
)

// abstract function to attach customizable event listener
type AbstractDecorator func(cards.AbstractGamestate) cards.AbstractGamestate

type ProgressListener struct {
	state            cards.AbstractGamestate
	MinibossDefeated int
}

func (s *ProgressListener) DoAction(data map[string]interface{}) {
	s.MinibossDefeated++
	if s.MinibossDefeated == 1 {
		// add in stronger hero in center deck
		deck := []cards.Card{}
		count := mrand.Int() % 2
		for i := 0; i < count; i++ {
			h := cards.NewMonsterSlayer(s.state)
			deck = append(deck, &h)

		}
		count = mrand.Int() % 2
		for i := 0; i < count; i++ {
			h := cards.NewStagShaman(s.state)
			deck = append(deck, &h)
		}
		count = mrand.Int() % 3
		for i := 0; i < count; i++ {
			h := cards.NewDeadweight(s.state)
			deck = append(deck, &h)
		}
		for i := 0; i < 3; i++ {
			h := cards.NewThief(s.state)
			deck = append(deck, &h)
		}
		// count = mrand.Int() % 5
		for i := 0; i < 4; i++ {
			h := cards.NewCleric(s.state)
			deck = append(deck, &h)
		}
		count = mrand.Int() % 5
		for i := 0; i < count; i++ {
			h := cards.NewShieldBasher(s.state)
			deck = append(deck, &h)
		}
		count = mrand.Int() % 3
		for i := 0; i < count; i++ {
			h := cards.NewArcher(s.state)
			deck = append(deck, &h)
		}
		count = mrand.Int() % 5
		for i := 0; i < count; i++ {
			h := cards.NewGoblinWolfRaiderMonster(s.state)
			deck = append(deck, &h)
		}
		for i := 0; i < 3; i++ {
			h := cards.NewNobleKnight(s.state)
			deck = append(deck, &h)
		}
		for i := 0; i < 2; i++ {
			h := cards.NewIceWyvern(s.state)
			deck = append(deck, &h)
		}
		for i := 0; i < 4; i++ {
			h := cards.NewTorchtail(s.state)
			deck = append(deck, &h)
		}
		for i := 0; i < 2; i++ {
			h := cards.NewSlimeLarge(s.state)
			deck = append(deck, &h)
		}
		for i := 0; i < 2; i++ {
			h := cards.NewFirelake(s.state)
			deck = append(deck, &h)
		}
		s.state.AddCardToCenterDeck(cards.DISCARD_SOURCE_NAN, true, deck...)
		// add more difficult location to explore

		return
	} else if s.MinibossDefeated == 2 {
		deck := []cards.Card{}
		for i := 0; i < 3; i++ {
			ll := cards.NewBulwark(s.state)
			deck = append(deck, &ll)
		}
		for i := 0; i < 4; i++ {
			h := cards.NewIceWyvern(s.state)
			deck = append(deck, &h)
		}
		s.state.AddCardToCenterDeck(cards.DISCARD_SOURCE_NAN, true, deck...)
	}
}

func AttachProgressionCounter(state cards.AbstractGamestate) cards.AbstractGamestate {
	state.AttachListener(cards.EVENT_MINIBOSS_DEFEATED, &ProgressListener{state: state})
	return state
}

// sample implementation of AbstractDecorator. This will add TombOfForgottenMonarch when the player recruit 3 Heroes
func AttachTombOfForgottenMonarch(state cards.AbstractGamestate) cards.AbstractGamestate {
	tomb := cards.NewTombMonarchEntrance(state)
	cardsAdded := []cards.Card{&tomb}
	pushCenterDeckAction := cards.NewPushCenterDeckAction(state, cardsAdded, false)
	removeEventListenerAction := cards.NewRemoveEventListenerAction(state, cards.EVENT_ATTR_CARD_RECRUITED, nil)
	compositeAction := cards.NewCompositeAction(state, pushCenterDeckAction, removeEventListenerAction)
	countDownAction := cards.NewCountDownAction(3, 1, compositeAction)
	addTombListener := cards.NewCardRecruitedListener(nil, countDownAction)
	removeEventListenerAction.(*cards.RemoveEventListenerAction).SetListener(addTombListener)

	state.AttachListener(cards.EVENT_CARD_RECRUITED, addTombListener)
	return state
}
