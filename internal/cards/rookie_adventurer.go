package cards

type RookieAdventurer struct {
	BaseHero
	gamestate AbstractGamestate
}

func NewRookieAdventurer(state AbstractGamestate) RookieAdventurer {
	this := RookieAdventurer{BaseHero: BaseHero{}}
	this.gamestate = state
	return this
}
func (r *RookieAdventurer) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
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
