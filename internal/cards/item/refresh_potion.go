package item

import (
	"github/kharism/GuildSim_go/internal/cards"
)

type RefreshPotion struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func NewRefreshPotion(state cards.AbstractGamestate) RefreshPotion {
	return RefreshPotion{state: state}
}

func (h *RefreshPotion) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func (h *RefreshPotion) GetName() string {
	return "RefreshPotion"
}
func (h *RefreshPotion) GetDescription() string {
	return "Discard all your hand then draw that many amount"
}
func (h *RefreshPotion) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

func (h *RefreshPotion) OnConsume() {
	cardsInHand := h.state.GetCardInHand()
	discardedCount := len(cardsInHand)
	for i := len(cardsInHand) - 1; i >= 0; i-- {
		curCard := cardsInHand[i]
		h.state.RemoveCardFromHandIdx(i)

		h.state.DiscardCard(curCard, cards.DISCARD_SOURCE_HAND)
	}
	for i := 0; i < discardedCount; i++ {
		h.state.Draw()
	}
	//h.state.AddResource(cards.RESOURCE_NAME_EXPLORATION, 3)
	h.state.RemoveItem(h)
	h.Dispose(cards.DISCARD_SOURCE_NAN)
}
