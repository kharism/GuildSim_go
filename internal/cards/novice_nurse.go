package cards

type NoviceNurse struct {
	BaseHero
	gamestate AbstractGamestate
}

func NewNoviceNurse(state AbstractGamestate) NoviceNurse {
	this := NoviceNurse{BaseHero: BaseHero{}}
	this.gamestate = state
	return this
}
func (r *NoviceNurse) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func (r *NoviceNurse) GetName() string {
	return "NoviceNurse"
}
func (r *NoviceNurse) GetCost() Cost {
	j := NewCost()
	j.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return j
}
func (r *NoviceNurse) GetDescription() string {
	return "Heal 2 HP, also yotu can add 1 card from your discard pile to top of the deck"
}
func (r *NoviceNurse) OnPlay() {
	r.gamestate.TakeDamage(-2)
	discardedCard := r.gamestate.GetCooldownCard()
	if len(discardedCard) > 0 {
		selectedCardId := r.gamestate.GetCardPicker().PickCardOptional(discardedCard, "Pick card to stacked to deck")
		returnedCard := discardedCard[selectedCardId]
		r.gamestate.RemoveCardFromCooldownIdx(selectedCardId)
		r.gamestate.StackCards(DISCARD_SOURCE_COOLDOWN, returnedCard)
	}
}
