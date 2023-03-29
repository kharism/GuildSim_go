package item

import "github/kharism/GuildSim_go/internal/cards"

type RecursionPotion struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func NewRecursionPotion(state cards.AbstractGamestate) RecursionPotion {
	return RecursionPotion{state: state}
}

func (h *RecursionPotion) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func (h *RecursionPotion) GetName() string {
	return "Recursion Potion"
}
func (h *RecursionPotion) GetDescription() string {
	return "add 1 card from cooldown pile to the top of the deck"
}
func (h *RecursionPotion) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

func (h *RecursionPotion) OnConsume() {
	cooldownPile := h.state.GetCooldownCard()
	idx := h.state.GetCardPicker().PickCard(cooldownPile, "select 1 to return to top of your deck")
	hh := cooldownPile[idx]
	h.state.RemoveCardFromCooldownIdx(idx)
	h.state.StackCards(cards.DISCARD_SOURCE_COOLDOWN, hh)
	//h.state.TakeDamage(-5)
	h.state.RemoveItem(h)
	h.Dispose(cards.DISCARD_SOURCE_NAN)
}
