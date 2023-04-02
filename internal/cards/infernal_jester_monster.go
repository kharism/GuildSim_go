package cards

type InfernalJester struct {
	BaseMonster
	gamestate AbstractGamestate
}

func (r *InfernalJester) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func NewInfernalJester(gamestate AbstractGamestate) InfernalJester {
	return InfernalJester{gamestate: gamestate}
}

func (r *InfernalJester) GetName() string {
	return "InfernalJester"
}
func (r *InfernalJester) GetDescription() string {
	return "recruitable. on punish: 1 damage. onslain: gain 6 reputation. on play: stack 1 card from your hand, then gain 4 combat point"
}
func (r *InfernalJester) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 3)
	return cost
}
func (r *InfernalJester) OnSlain() {
	r.gamestate.AddResource(RESOURCE_NAME_REPUTATION, 6)
}
func (r *InfernalJester) OnPlay() {
	cardInHand := r.gamestate.GetCardInHand()
	if len(cardInHand) > 0 {
		selection := r.gamestate.GetCardPicker().PickCard(cardInHand, "Pick 1 to stack to deck")
		selectedCard := cardInHand[selection]
		r.gamestate.RemoveCardFromHandIdx(selection)
		r.gamestate.StackCards(DISCARD_SOURCE_HAND, selectedCard)
		r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 4)
	}

}
