package cards

type GoblinRaiderMonster struct {
	BaseMonster
	state AbstractGamestate
}

func NewGoblinWolfRaiderMonster(state AbstractGamestate) GoblinRaiderMonster {
	return GoblinRaiderMonster{state: state}

}
func (m *GoblinRaiderMonster) GetName() string {
	return "GoblinRaiderMonster"
}
func (m *GoblinRaiderMonster) GetDescription() string {
	return "2 damage per turn"
}
func (m *GoblinRaiderMonster) GetCost() Cost {
	cost := NewCost()
	cost.Detail[RESOURCE_NAME_COMBAT] = 3
	return cost
}
func (m *GoblinRaiderMonster) OnPunish() {
	m.state.TakeDamage(2)
}
func (m *GoblinRaiderMonster) Dispose(source string) {
	m.state.DiscardCard(m, source)
}
