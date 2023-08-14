package item

import (
	"github/kharism/GuildSim_go/internal/cards"
)

type CombatGauntlet struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func NewCombatGauntlet(state cards.AbstractGamestate) CombatGauntlet {
	return CombatGauntlet{state: state}
}

func (h *CombatGauntlet) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func (h *CombatGauntlet) GetName() string {
	return "CombatGauntlet"
}
func (h *CombatGauntlet) GetDescription() string {
	return "gain 2 combat point each turn"
}
func (h *CombatGauntlet) OnAcquire() {
	action := cards.NewAddResourceAction(h.state, cards.RESOURCE_NAME_COMBAT, 2)
	listener := cards.NewBasicAction(action)
	h.state.AttachListener(cards.EVENT_START_OF_TURN, listener)
}
func (h *CombatGauntlet) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}
