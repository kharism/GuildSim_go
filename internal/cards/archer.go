package cards

type Archer struct {
	BaseHero
	gamestate AbstractGamestate
}

func NewArcher(state AbstractGamestate) Archer {
	h := Archer{}
	h.gamestate = state
	return h
}
func (h *Archer) GetName() string {
	return "Archer"
}
func (h *Archer) GetDescription() string {
	return "gain combat equal to your block"
}
func (h *Archer) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}
func (h *Archer) Dispose(source string) {
	h.gamestate.DiscardCard(h, source)
}
func (h *Archer) OnPlay() {
	playedCards := h.gamestate.GetPlayedCards()
	h.gamestate.AddResource(RESOURCE_NAME_COMBAT, len(playedCards))
}
