package cards

type SlimeSmall struct {
	BaseMonster
	state AbstractGamestate
	// turnCounter int
}

func NewSlimeSmall(state AbstractGamestate) SlimeSmall {
	k := SlimeSmall{state: state}
	return k
}
func (m *SlimeSmall) GetName() string {
	return "Slime(S)"
}
func (m *SlimeSmall) GetDescription() string {
	return "deals 1 damage. on defeat, gains 1 reputation"
}
func (m *SlimeSmall) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 1)
	return cost
}
func (m *SlimeSmall) OnPunish() {
	m.state.TakeDamage(2)
}
func (m *SlimeSmall) OnSlain() {
	// m.state.AddResource(RESOURCE_NAME_EXPLORATION, 1)
	m.state.AddResource(RESOURCE_NAME_REPUTATION, 2)
	// h := []Card{}
	// for i := 0; i < 3; i++ {
	// 	midSlime := NewSlimeMid(m.state)
	// 	h = append(h, &midSlime)
	// }
	// m.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, false, h...)
}
