package item

import (
	"github/kharism/GuildSim_go/internal/cards"
)

type HealingPotion struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func (h *HealingPotion) Dispose(source string) {

}
func (h *HealingPotion) GetName() string {
	return "Healing Potion"
}
func (h *HealingPotion) GetDescription() string {
	return "Heal 5 damage"
}
func (h *HealingPotion) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

func (h *HealingPotion) OnConsume(source string) {
	h.state.TakeDamage(-5)
}
