package cards

type GoblinWolfRaiderMonster struct {
	BaseMonster
	state AbstractGamestate
}

func NewGoblinWolfRaiderMonster(state AbstractGamestate) GoblinWolfRaiderMonster {
	return GoblinWolfRaiderMonster{state: state}

}
func (m *GoblinWolfRaiderMonster) GetName() string {
	return "GoblinWolfRaiderMonster"
}
func (m *GoblinWolfRaiderMonster) GetDescription() string {
	return "GoblinWolfRaiderMonster"
}
func (m *GoblinWolfRaiderMonster) GetCost() Cost {
	cost := NewCost()
	cost.Detail[RESOURCE_NAME_COMBAT] = 3
	return cost
}
func (m *GoblinWolfRaiderMonster) OnPunish() {
	m.state.TakeDamage(2)
}
func (m *GoblinWolfRaiderMonster) Dispose(source string) {
	m.state.DiscardCard(m, source)
}
