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
func (r *RookieNurse) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func (r *RookieNurse) GetName() string {
	return "RookieNurse"
}
func (r *RookieNurse) GetDescription() string {
	return "Heal 2 Hp"
}
func (r *RookieNurse) GetCost() Cost {
	cost := NewCost()
	cost.Detail[RESOURCE_NAME_EXPLORATION] = 1
	return cost
}
func (r *RookieNurse) OnPlay() {
	// r.gamestate.Draw()
	r.gamestate.TakeDamage(-1)
}
