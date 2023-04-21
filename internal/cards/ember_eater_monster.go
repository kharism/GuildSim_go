package cards

type EmberEater struct {
	BaseMonster
	gamestate AbstractGamestate
}

func NewEmberEater(gamestate AbstractGamestate) EmberEater {
	return EmberEater{gamestate: gamestate}
}
func (r *EmberEater) GetName() string {
	return "EmberEater"
}
func (r *EmberEater) GetDescription() string {
	return "on punish: 3 damage. onslain: gain 3 reputation and draw 1 card"
}
func (r *EmberEater) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 4)
	return cost
}
func (m *EmberEater) OnPunish() {
	m.gamestate.TakeDamage(3)
}
func (r *EmberEater) OnSlain() {
	r.gamestate.AddResource(RESOURCE_NAME_REPUTATION, 3)
	r.gamestate.Draw()
}
