package cards

type RookieHunter struct {
	BaseHero
	state AbstractGamestate
}

func (r *RookieHunter) Dispose(source string) {
	r.state.DiscardCard(r, source)
}

func NewRookieHunter(gamestate AbstractGamestate) RookieHunter {
	return RookieHunter{state: gamestate}
}

func (r *RookieHunter) GetName() string {
	return "RookieHunter"
}
func (r *RookieHunter) GetDescription() string {
	return "gain either 2 Exploration or 2 Combat"
}
func (r *RookieHunter) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}
func (r *RookieHunter) OnPlay() {
	if r.state.GetBoolPicker().BoolPick("Gain 2 exploration?") {
		r.state.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	} else {
		//r.gamestate.AddResource(RESOURCE_NAME_BLOCK, 5)
		r.state.AddResource(RESOURCE_NAME_COMBAT, 2)
	}

}
