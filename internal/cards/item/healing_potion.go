package item

import (
	"github/kharism/GuildSim_go/internal/cards"
)

type HealingPotion struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

<<<<<<< HEAD
func NewHealingPotion(state cards.AbstractGamestate) HealingPotion {
	return HealingPotion{state: state}
}

func (h *HealingPotion) Dispose(source string) {
	h.state.BanishCard(h, source)
=======
func (h *HealingPotion) Dispose(source string) {

>>>>>>> 3ea1550 (add item related stuff.)
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

<<<<<<< HEAD
func (h *HealingPotion) OnConsume() {
	h.state.TakeDamage(-5)
	h.state.RemoveItem(h)
	h.Dispose(cards.DISCARD_SOURCE_NAN)
=======
func (h *HealingPotion) OnConsume(source string) {
	h.state.TakeDamage(-5)
>>>>>>> 3ea1550 (add item related stuff.)
}
