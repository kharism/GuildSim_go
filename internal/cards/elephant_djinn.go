package cards

type ElephantDjinn struct {
	BaseMonster
	gamestate AbstractGamestate
}

func (b *ElephantDjinn) GetName() string {
	return "ElephantDjinn"
}
func (b *ElephantDjinn) GetDescription() string {
	return "recruitable. on recruit draw 3 cards then discard 1. on play: get 4 block"
}
func (b *ElephantDjinn) Dispose(source string) {
	b.gamestate.DiscardCard(b, source)
}
func NewElephantDjinn(state AbstractGamestate) ElephantDjinn {
	b := ElephantDjinn{gamestate: state}
	return b
}
func (b *ElephantDjinn) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 2)
	return cost
}
func (r *ElephantDjinn) OnSlain() {
	r.gamestate.AddResource(RESOURCE_NAME_REPUTATION, 3)
}
func (r *ElephantDjinn) OnRecruit() {
	r.gamestate.Draw()
	r.gamestate.Draw()
	r.gamestate.Draw()
	hand := r.gamestate.GetCardInHand()
	idx := r.gamestate.GetCardPicker().PickCard(hand, "Discard 1")
	discardedCard := hand[idx]
	r.gamestate.DiscardCard(discardedCard, DISCARD_SOURCE_HAND)
}
func (r *ElephantDjinn) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_BLOCK, 4)
}
