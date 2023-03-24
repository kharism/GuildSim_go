package cards

import "fmt"

type PushCenterDeckAction struct {
	state      AbstractGamestate
	shuffle    bool
	cardsAdded []Card
}

func NewPushCenterDeckAction(state AbstractGamestate, cards []Card, shuffle bool) *PushCenterDeckAction {
	return &PushCenterDeckAction{state: state, cardsAdded: cards, shuffle: shuffle}
}
func (p *PushCenterDeckAction) DoAction() {
	fmt.Println("Add card to center deck", p.cardsAdded[0].GetName())
	p.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, p.shuffle, p.cardsAdded...)
}
