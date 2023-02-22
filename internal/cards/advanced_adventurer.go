package cards

type AdvancedAdventurer struct {
	RookieAdventurer
}

func NewAdvancedAdventurer(state AbstractGamestate) AdvancedAdventurer {
	base := NewRookieAdventurer(state)
	this := AdvancedAdventurer{base}
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
