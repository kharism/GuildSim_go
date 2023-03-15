package cards

type PushCenterDeckAction struct {
	state      AbstractGamestate
	cardsAdded []Card
}

func NewPushCenterDeckAction(state AbstractGamestate, cards []Card) *PushCenterDeckAction {
	return &PushCenterDeckAction{state: state, cardsAdded: cards}
}
func (p *PushCenterDeckAction) DoAction() {
	p.state.AddCardToCenterDeck(p.cardsAdded...)
}
