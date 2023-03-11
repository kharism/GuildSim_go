package cards

type AdvancedAdventurer struct {
	BaseHero
	gamestate AbstractGamestate
}

func NewAdvancedAdventurer(state AbstractGamestate) AdvancedAdventurer {
	// base := NewRookieAdventurer(state)
	this := AdvancedAdventurer{gamestate: state}
	return this
}

func (r *AdvancedAdventurer) GetName() string {
	return "AdvancedAdventurer"
}
func (r *AdvancedAdventurer) GetDescription() string {
	return "Add 2 exlporation point"
}
func (r *AdvancedAdventurer) GetCost() Cost {
	cost := NewCost()
	cost.Resource.Detail[RESOURCE_NAME_MONEY] = 100
	return cost
}
func (r *AdvancedAdventurer) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 2)
}
func (r *AdvancedAdventurer) Dispose() {
	r.gamestate.DiscardCard(r)
}
