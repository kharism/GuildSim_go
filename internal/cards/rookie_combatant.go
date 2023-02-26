package cards

type RookieCombatant struct {
	BaseHero
	gamestate AbstractGamestate
}

func NewRookieCombatant(state AbstractGamestate) RookieCombatant {
	this := RookieCombatant{BaseHero: BaseHero{}}
	this.gamestate = state
	return this
}

func (r *RookieCombatant) GetName() string {
	return "RookieCombatant"
}
func (r *RookieCombatant) GetDescription() string {
	return "Add 1 Combat point"
}
func (r *RookieCombatant) GetCost() Cost {
	cost := NewCost()
	return cost
}
func (r *RookieCombatant) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 1)
}
