package cards

type GoblinMinionMonster struct {
	BaseMonster
	state AbstractGamestate
}

func NewGoblinMinionMonster(state AbstractGamestate) GoblinMinionMonster {
	g := GoblinMinionMonster{state: state}
	return g
}

func (ed *GoblinMinionMonster) GetName() string {
	return "goblinMinion"
}
func (ed *GoblinMinionMonster) Dispose(source string) {
	ed.state.BanishCard(ed, source)
}
func (ed *GoblinMinionMonster) GetDescription() string {
	return "2 dmg per turn"
}
func (ed *GoblinMinionMonster) GetCost() Cost {
	c := NewCost()
	c.AddResource(RESOURCE_NAME_COMBAT, 2)
	return c
}
func (g *GoblinMinionMonster) OnPunish() {
	g.state.TakeDamage(2)
}
func (g *GoblinMinionMonster) OnSlain() {
	//g.state.AddResource(RESOURCE_NAME_REPUTATION, 1)
}
