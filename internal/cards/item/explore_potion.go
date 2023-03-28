package item

import "github/kharism/GuildSim_go/internal/cards"

type ExplorePotion struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func NewExplorePotion(state cards.AbstractGamestate) ExplorePotion {
	return ExplorePotion{state: state}
}

func (h *ExplorePotion) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func (h *ExplorePotion) GetName() string {
	return "Explore Potion"
}
func (h *ExplorePotion) GetDescription() string {
	return "Add 3 Explore"
}
func (h *ExplorePotion) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

func (h *ExplorePotion) OnConsume() {
	h.state.AddResource(cards.RESOURCE_NAME_EXPLORATION, 3)
	h.state.RemoveItem(h)
	h.Dispose(cards.DISCARD_SOURCE_NAN)
}
