package cards

type PyroKnight struct {
	BaseHero
	gamestate AbstractGamestate
}

func (r *PyroKnight) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func NewPyroKnight(gamestate AbstractGamestate) PyroKnight {
	return PyroKnight{gamestate: gamestate}
}

func (r *PyroKnight) GetName() string {
	return "Pyro Knight"
}
func (r *PyroKnight) GetDescription() string {
	return "Add 1 Combat point and discard 1 card"
}
func (r *PyroKnight) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}
func (r *PyroKnight) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 2)
	cardPicker := r.gamestate.GetCardPicker()
	cardInHand := r.gamestate.GetCardInHand()
	if len(cardInHand) > 0 {
		discardedIdx := cardPicker.PickCard(cardInHand, "Choose Card to discard")

		discardedCard := cardInHand[discardedIdx]
		r.gamestate.RemoveCardFromHandIdx(discardedIdx)

		r.gamestate.DiscardCard(discardedCard, DISCARD_SOURCE_HAND)
	}
}
