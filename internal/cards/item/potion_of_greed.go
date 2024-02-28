package item

import "github/kharism/GuildSim_go/internal/cards"

type GreedPotion struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func NewGreedPotion(state cards.AbstractGamestate) GreedPotion {
	return GreedPotion{state: state}
}

func (h *GreedPotion) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func (h *GreedPotion) GetName() string {
	return "GreedPotion"
}
func (h *GreedPotion) GetDescription() string {
	return "draw 2 cards"
}
func (h *GreedPotion) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

func (h *GreedPotion) OnConsume() {
	//h.state.TakeDamage(-5)
	h.state.Draw()
	h.state.Draw()
	h.state.RemoveItem(h)
	h.Dispose(cards.DISCARD_SOURCE_NAN)
}
