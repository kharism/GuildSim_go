package cards

type GoblinMonster struct {
	BaseMonster
	state AbstractGamestate
}

func NewGoblinMonster(state AbstractGamestate) GoblinMonster {
	g := GoblinMonster{state: state}
	return g
}

func (ed *GoblinMonster) GetName() string {
	return "GoblinMonster"
}
func (ed *GoblinMonster) Dispose(source string) {
	ed.state.BanishCard(ed, source)
}
func (ed *GoblinMonster) GetDescription() string {
	return "1 dmg per turn, Reward: 1 Reputation"
}
func (ed *GoblinMonster) GetCost() Cost {
	c := NewCost()
	c.AddResource(RESOURCE_NAME_COMBAT, 1)
	return c
}
func (g *GoblinMonster) OnPunish() {
	g.state.TakeDamage(1)
}
func (g *GoblinMonster) OnSlain() {
	g.state.AddResource(RESOURCE_NAME_REPUTATION, 1)
}
