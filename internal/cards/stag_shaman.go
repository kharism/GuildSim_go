package cards

type StagShaman struct {
	BaseHero
	state AbstractGamestate
}

func NewStagShaman(state AbstractGamestate) StagShaman {
	s := StagShaman{state: state}
	return s
}
func (s *StagShaman) GetName() string {
	return "Stag Shaman"
}
func (h *StagShaman) GetDescription() string {
	return "You can banish 1 card from Cooldown pile then draw 1 card regardless"
}
func (h *StagShaman) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return cost
}

func (h *StagShaman) Dispose() {
	h.state.DiscardCard(h)
}

func (h *StagShaman) OnPlay() {
	cooldownList := h.state.GetCooldownCard()
	if len(cooldownList) > 0 {
		cardPicker := h.state.GetCardPicker()
		cardId := cardPicker.PickCardOptional(cooldownList, "Pick a card to banish")
		card := cooldownList[cardId]
		h.state.RemoveCardFromCooldownIdx(cardId)
		h.state.BanishCard(card)
	}
	h.state.Draw()

}
