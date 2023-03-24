package cards

type RookieMage struct {
	BaseHero
	state AbstractGamestate
}

func NewRookieMage(state AbstractGamestate) RookieMage {
	j := RookieMage{state: state}
	return j
}

func (h *RookieMage) GetName() string {
	return "Rookie Mage"
}
func (h *RookieMage) GetDescription() string {
	return "Discard 1 then draw 2"
}
func (h *RookieMage) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}

func (h *RookieMage) OnPlay() {
	cardPicker := h.state.GetCardPicker()
	cardInHand := h.state.GetCardInHand()
	if len(cardInHand) > 0 {
		discardedIdx := cardPicker.PickCard(cardInHand, "Choose Card to discard")

		discardedCard := cardInHand[discardedIdx]

		h.state.RemoveCardFromHandIdx(discardedIdx)

		h.state.DiscardCard(discardedCard, DISCARD_SOURCE_HAND)
		for i := 0; i < 2; i++ {
			h.state.Draw()
		}
	}

}

func (h *RookieMage) Dispose(source string) {
	h.state.DiscardCard(h, source)
}
