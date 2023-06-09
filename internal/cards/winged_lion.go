package cards

import "fmt"

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
	shuffleableCard := []Card{}
	for _, c := range cardList {
		if _, ok := c.(Unshuffleable); ok {
			continue
		} else {
			shuffleableCard = append(shuffleableCard, c)
		}
	}
	cardPicker := h.state.GetCardPicker()

	// idx_shuf is not always the same with index of cards in center due to unshuffleable cards
	idx_shuf := cardPicker.PickCard(shuffleableCard, "Pick a card to shuffle to deck then draw")
	// newCenterCards := []Card{}
	c := shuffleableCard[idx_shuf]

	real_idx := -1
	for idx, v := range h.state.GetCenterCard() {
		if v == c {
			//newCenterCards = append(newCenterCards, replacementCard)
			real_idx = idx
		} else {
			//newCenterCards = append(newCenterCards, v)
		}
	}
	h.state.AddCardToCenterDeck(DISCARD_SOURCE_CENTER, false, c)
	h.state.RemoveCardFromCenterRowIdx(real_idx)
	replacementCard := h.state.ReplaceCenterCard()
	fmt.Println("Replace", c.GetName(), "with", replacementCard.GetName())
	h.state.AppendCenterCard(replacementCard)
	// newCenterCards = append(newCenterCards, replacementCard)
	// h.state.CenterCards = newCenterCards
	// selectedCard := cardList[idx]
	// h.state.RemoveCardFromCenterRowIdx(idx)
	// topdeck := h.state.ReplaceCenterCard()
	// h.state.AppendCenterCard(topdeck)
	// // cardList[idx] = topdeck
	// h.state.AddCardToCenterDeck(DISCARD_SOURCE_CENTER, false, selectedCard)
	h.state.Draw()
}

func (h *WingedLion) Dispose(source string) {
	h.state.DiscardCard(h, source)
}
