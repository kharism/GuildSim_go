package item

import "github/kharism/GuildSim_go/internal/cards"

type CombatPotion struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func NewCombatPotion(state cards.AbstractGamestate) CombatPotion {
	return CombatPotion{state: state}
}

func (h *CombatPotion) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func (h *CombatPotion) GetName() string {
	return "Combat Potion"
}
func (h *CombatPotion) GetDescription() string {
	return "Add 3 Combat"
}
func (h *CombatPotion) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

func (h *CombatPotion) OnConsume() {
	h.state.AddResource(cards.RESOURCE_NAME_COMBAT, 3)
	h.state.RemoveItem(h)
	h.Dispose(cards.DISCARD_SOURCE_NAN)
}
