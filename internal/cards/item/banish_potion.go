package item

import "github/kharism/GuildSim_go/internal/cards"

type BanishPotion struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func NewBanishPotion(state cards.AbstractGamestate) BanishPotion {
	return BanishPotion{state: state}
}

func (h *BanishPotion) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func (h *BanishPotion) GetName() string {
	return "BanishPotion"
}
func (h *BanishPotion) GetDescription() string {
	return "Banish 1 card from discard pile"
}
func (h *BanishPotion) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

func (h *BanishPotion) OnConsume() {
	// h.state.AddResource(cards.RESOURCE_NAME_COMBAT, 3)
	discardPile := h.state.GetCooldownCard()
	pickedCardIdx := h.state.GetCardPicker().PickCard(discardPile, "Banish 1 card from cooldown pile")
	pickedCard := discardPile[pickedCardIdx]
	h.state.RemoveCardFromCooldownIdx(pickedCardIdx)
	h.state.BanishCard(pickedCard, cards.DISCARD_SOURCE_COOLDOWN)
	h.state.RemoveItem(h)
	h.Dispose(cards.DISCARD_SOURCE_NAN)
}
