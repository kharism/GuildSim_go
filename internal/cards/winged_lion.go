package cards

type WingedLion struct {
	BaseHero
	state AbstractGamestate
}

func NewWingedLion(state AbstractGamestate) WingedLion {
	j := WingedLion{state: state}
	return j
}

func (h *WingedLion) GetName() string {
	return "Winged Lion"
}
func (h *WingedLion) GetDescription() string {
	return "replace 1 card from center row with top deck then draw"
}
func (h *WingedLion) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}
func (h *WingedLion) OnPlay() {
	cardList := h.state.GetCenterCard()
	cardPicker := h.state.GetCardPicker()

	idx := cardPicker.PickCard(cardList, "Pick a card to shuffle to deck then draw")
	selectedCard := cardList[idx]
	topdeck := h.state.ReplaceCenterCard()
	cardList[idx] = topdeck
	h.state.AddCardToCenterDeck(selectedCard)
	h.state.Draw()
}

func (h *WingedLion) Dispose(source string) {
	h.state.DiscardCard(h, source)
}
