package item

import (
	"github/kharism/GuildSim_go/internal/cards"
)

type CombatTalisman struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func NewCombatTalisman(state cards.AbstractGamestate) CombatTalisman {
	return CombatTalisman{state: state}
}

func (h *CombatTalisman) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func (h *CombatTalisman) GetName() string {
	return "CombatTalisman"
}
func (h *CombatTalisman) GetDescription() string {
	return "gain 2 combat point each turn"
}
func (h *CombatTalisman) OnAcquire() {
	action := cards.NewAddResourceAction(h.state, cards.RESOURCE_NAME_COMBAT, 2)
	listener := cards.NewBasicAction(action)
	h.state.AttachListener(cards.EVENT_START_OF_TURN, listener)
}
func (h *CombatTalisman) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}
