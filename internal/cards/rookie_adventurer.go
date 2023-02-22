package cards

type RookieAdventurer struct {
	gamestate AbstractGamestate
}

func NewRookieAdventurer(state AbstractGamestate) RookieAdventurer {
	this := RookieAdventurer{}
	this.gamestate = state
	return this
}

func (r *RookieAdventurer) GetName() string {
	return "RookieAdventurer"
}
func (r *RookieAdventurer) GetDescription() string {
	return "Add 1 exlporation point"
}
func (r *RookieAdventurer) GetCost() Cost {
	cost := NewCost()
	return cost
}
func (r *RookieAdventurer) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 1)
}
