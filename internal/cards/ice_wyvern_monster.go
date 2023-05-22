package cards

type IceWyvern struct {
	BaseMonster
	state AbstractGamestate
}

func NewIceWyvern(state AbstractGamestate) IceWyvern {
	k := IceWyvern{state: state}
	return k
}
func (m *IceWyvern) GetName() string {
	return "IceWyvern"
}
func (m *IceWyvern) GetDescription() string {
	return "deals 1 damage, Also stack 1 freeze curse. On defeat: gains 2 reputation and draw 2 cards"
}
func (m *IceWyvern) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 4)
	return cost
}
func (m *IceWyvern) OnPunish() {
	m.state.TakeDamage(1)
	curse1 := NewFreezeCurse(m.state)
	m.state.StackCards(DISCARD_SOURCE_NAN, &curse1)
}
func (m *IceWyvern) OnSlain() {
	// m.state.AddResource(RESOURCE_NAME_EXPLORATION, 1)
	m.state.AddResource(RESOURCE_NAME_REPUTATION, 2)
	// h := []Card{}
	for i := 0; i < 2; i++ {
		// midSlime := NewSlimeMid(m.state)
		// h = append(h, &midSlime)
		m.state.Draw()
	}
	// m.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, false, h...)
}
