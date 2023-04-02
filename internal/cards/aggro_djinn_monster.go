package cards

type AggroDjinn struct {
	BaseMonster
	gamestate AbstractGamestate
}

func (r *AggroDjinn) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func NewAggroDjinn(gamestate AbstractGamestate) AggroDjinn {
	return AggroDjinn{gamestate: gamestate}
}
func (r *AggroDjinn) OnRecruit() {

}
func (r *AggroDjinn) GetName() string {
	return "AggroDjinn"
}
func (r *AggroDjinn) GetDescription() string {
	return "recruitable. onslain: gain 3 reputation. on play: gain 3 combat then add a copy of AggroDjinn to cooldown pile"
}
func (r *AggroDjinn) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 3)
	return cost
}
func (r *AggroDjinn) OnSlain() {
	r.gamestate.AddResource(RESOURCE_NAME_REPUTATION, 3)
}
func (r *AggroDjinn) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 3)
	newAggroDjinn := NewAggroDjinn(r.gamestate)
	r.gamestate.DiscardCard(&newAggroDjinn, DISCARD_SOURCE_NAN)

}
