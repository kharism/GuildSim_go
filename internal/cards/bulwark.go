package cards

type Bulwark struct {
	BaseHero
	gamestate AbstractGamestate
}

func (r *Bulwark) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func NewBulwark(gamestate AbstractGamestate) Bulwark {
	return Bulwark{gamestate: gamestate}
}

func (r *Bulwark) GetName() string {
	return "Bulwark"
}
func (r *Bulwark) GetDescription() string {
	return "gain 2 Combat or 5 block"
}
func (r *Bulwark) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return cost
}
func (r *Bulwark) OnPlay() {
	if r.gamestate.GetBoolPicker().BoolPick("Gain 2 combat?") {
		r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 2)
	} else {
		r.gamestate.AddResource(RESOURCE_NAME_BLOCK, 5)
	}

}
