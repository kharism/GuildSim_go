package cards

type RookieNurse struct {
	BaseHero
	gamestate AbstractGamestate
}

func NewRookieNurse(state AbstractGamestate) RookieNurse {
	n := RookieNurse{BaseHero: BaseHero{}}
	n.gamestate = state
	return n
}
func (r *RookieNurse) Dispose() {
	r.gamestate.DiscardCard(r)
}
func (r *RookieNurse) GetName() string {
	return "RookieCombatant"
}
func (r *RookieNurse) GetDescription() string {
	return "Draw 1 card"
}
func (r *RookieNurse) GetCost() Cost {
	cost := NewCost()
	cost.Detail[RESOURCE_NAME_EXPLORATION] = 1
	return cost
}
func (r *RookieNurse) OnPlay() {
	r.gamestate.Draw()
}
