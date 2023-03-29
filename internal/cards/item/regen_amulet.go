package item

import "github/kharism/GuildSim_go/internal/cards"

type RegenAmulet struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func NewRegenAmulet(state cards.AbstractGamestate) RegenAmulet {
	return RegenAmulet{state: state}
}

func (h *RegenAmulet) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func (h *RegenAmulet) GetName() string {
	return "Healing Potion"
}
func (h *RegenAmulet) GetDescription() string {
	return "Heal 1 HP each turn"
}
func (h *RegenAmulet) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

type HealAction struct {
	state cards.AbstractGamestate
}

func (h *HealAction) DoAction() {
	h.state.TakeDamage(-1)
}
func (h *RegenAmulet) OnAcquire() {
	action := &HealAction{state: h.state} //cards.NewAddResourceAction(h.state, cards.RESOURCE_NAME_EXPLORATION, 2)
	listener := cards.NewBasicAction(action)
	h.state.AttachListener(cards.EVENT_START_OF_TURN, listener)
}
