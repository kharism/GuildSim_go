package item

import "github/kharism/GuildSim_go/internal/cards"

type CompanionBuckler struct {
	cards.BaseItem
	gamestate cards.AbstractGamestate
}

func NewCompanionBuckler(state cards.AbstractGamestate) CompanionBuckler {
	return CompanionBuckler{gamestate: state}
}

func (h *CompanionBuckler) Dispose(source string) {
	h.gamestate.BanishCard(h, source)
}
func (h *CompanionBuckler) GetName() string {
	return "CompanionBuckler"
}
func (h *CompanionBuckler) GetDescription() string {
	return "each time you recruit a hero, gain 3 block"
}
func (h *CompanionBuckler) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

type GainBlockAction struct {
	state cards.AbstractGamestate
}

func (h *GainBlockAction) DoAction() {
	h.state.AddResource(cards.RESOURCE_NAME_BLOCK, 3)
}

func (h *CompanionBuckler) OnAcquire() {
	action := &GainBlockAction{state: h.gamestate} //cards.NewAddResourceAction(h.state, cards.RESOURCE_NAME_EXPLORATION, 2)
	listener := cards.NewBasicAction(action)
	h.gamestate.AttachListener(cards.EVENT_CARD_RECRUITED, listener)
}
