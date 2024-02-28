package cards

type NoviceAdventurer struct {
	BaseHero
	gamestate AbstractGamestate
}

func NewNoviceAdventurer(state AbstractGamestate) NoviceAdventurer {
	this := NoviceAdventurer{BaseHero: BaseHero{}}
	this.gamestate = state
	return this
}
func (r *NoviceAdventurer) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func (r *NoviceAdventurer) GetName() string {
	return "NoviceAdventurer"
}
func (r *NoviceAdventurer) GetCost() Cost {
	j := NewCost()
	j.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return j
}
func (r *NoviceAdventurer) GetDescription() string {
	return "gain 3 exploration and 3 block"
}
func (r *NoviceAdventurer) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_BLOCK, 3)
	r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 3)
}
