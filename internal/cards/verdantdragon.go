package cards

type VerdantDragon struct {
	BaseHero
	state AbstractGamestate
}

func NewVerdantDragon(state AbstractGamestate) VerdantDragon {
	return VerdantDragon{state: state}
}

func (v *VerdantDragon) GetName() string {
	return "VerdantDragon"
}
func (v *VerdantDragon) GetDescription() string {
	return "on play: stack 1 card to main deck, then gain 7HP and 2 block"
}
func (v *VerdantDragon) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 5)
	return cost
}
func (v *VerdantDragon) Dispose(source string) {
	v.state.DiscardCard(v, source)
}
func (v *VerdantDragon) OnPlay() {
	cardInHand := v.state.GetCardInHand()
	if len(cardInHand) > 0 {
		selection := v.state.GetCardPicker().PickCard(cardInHand, "Pick 1 to stack to deck")
		selectedCard := cardInHand[selection]
		v.state.RemoveCardFromHandIdx(selection)
		v.state.StackCards(DISCARD_SOURCE_HAND, selectedCard)
		v.state.TakeDamage(-7)
		v.state.AddResource(RESOURCE_NAME_BLOCK, 2)
	}
}
