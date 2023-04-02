package decorator

import "github/kharism/GuildSim_go/internal/cards"

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
		for i := 0; i < 2; i++ {
			h := cards.NewMonsterSlayer(s.state)
			deck = append(deck, &h)

		}
		for i := 0; i < 2; i++ {
			h := cards.NewStagShaman(s.state)
			deck = append(deck, &h)
		}
		for i := 0; i < 3; i++ {
			h := cards.NewThief(s.state)
			deck = append(deck, &h)
		}
		for i := 0; i < 3; i++ {
			h := cards.NewNobleKnight(s.state)
			deck = append(deck, &h)
		}
		for i := 0; i < 5; i++ {
			h := cards.NewSlimeLarge(s.state)
			deck = append(deck, &h)
		}
		s.state.AddCardToCenterDeck(cards.DISCARD_SOURCE_NAN, true, deck...)
		// add more difficult location to explore
		return
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
