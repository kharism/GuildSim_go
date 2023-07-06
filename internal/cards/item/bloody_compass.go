package item

import "github/kharism/GuildSim_go/internal/cards"

type BloodyCompass struct {
	cards.BaseItem
	gamestate cards.AbstractGamestate
}

func NewBloodyCompass(state cards.AbstractGamestate) BloodyCompass {
	return BloodyCompass{gamestate: state}
}

func (h *BloodyCompass) Dispose(source string) {
	h.gamestate.BanishCard(h, source)
}
func (h *BloodyCompass) GetName() string {
	return "BloodyCompass"
}
func (h *BloodyCompass) GetDescription() string {
	return "each time a monster is defeated, gains 2 exploration"
}
func (h *BloodyCompass) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

type ExploreMoreAction struct {
	state cards.AbstractGamestate
}

func (h *ExploreMoreAction) DoAction() {
	h.state.AddResource(cards.RESOURCE_NAME_EXPLORATION, 2)
}

func (h *BloodyCompass) OnAcquire() {
	action := &ExploreMoreAction{state: h.gamestate} //cards.NewAddResourceAction(h.state, cards.RESOURCE_NAME_EXPLORATION, 2)
	listener := cards.NewBasicAction(action)
	h.gamestate.AttachListener(cards.EVENT_CARD_DEFEATED, listener)
}
