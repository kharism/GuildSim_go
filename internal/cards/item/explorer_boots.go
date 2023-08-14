package item

import "github/kharism/GuildSim_go/internal/cards"

type ExplorerBoots struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func (h *ExplorerBoots) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func NewExplorerBoots(state cards.AbstractGamestate) ExplorerBoots {
	return ExplorerBoots{state: state}
}
func (h *ExplorerBoots) GetName() string {
	return "ExplorerBoots"
}
func (h *ExplorerBoots) GetDescription() string {
	return "gain 2 exploration point each turn"
}
func (h *ExplorerBoots) OnAcquire() {
	action := cards.NewAddResourceAction(h.state, cards.RESOURCE_NAME_EXPLORATION, 2)
	listener := cards.NewBasicAction(action)
	h.state.AttachListener(cards.EVENT_START_OF_TURN, listener)
}
func (h *ExplorerBoots) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}
