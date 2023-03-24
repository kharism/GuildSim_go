package decorator

import "github/kharism/GuildSim_go/internal/cards"

// abstract function to attach customizable event listener
type AbstractDecorator func(cards.AbstractGamestate) cards.AbstractGamestate

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
