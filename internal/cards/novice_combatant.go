package cards

type NoviceCombatant struct {
	BaseHero
	gamestate AbstractGamestate
}

func NewNoviceCombatant(state AbstractGamestate) NoviceCombatant {
	this := NoviceCombatant{BaseHero: BaseHero{}}
	this.gamestate = state
	return this
}
func (r *NoviceCombatant) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func (r *NoviceCombatant) GetName() string {
	return "NoviceCombatant"
}
func (r *NoviceCombatant) GetCost() Cost {
	j := NewCost()
	j.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return j
}
func (r *NoviceCombatant) GetDescription() string {
	return "gain 4 combat and 3 block"
}
func (r *NoviceCombatant) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_BLOCK, 3)
	r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 4)
}
