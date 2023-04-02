package cards

type SlimeMid struct {
	BaseMonster
	state       AbstractGamestate
	turnCounter int
}

func NewSlimeMid(state AbstractGamestate) SlimeMid {
	k := SlimeMid{state: state}
	return k
}
func (m *SlimeMid) GetName() string {
	return "Slime(M)"
}
func (m *SlimeMid) GetDescription() string {
	return "deals 2 damage.on defeat, gains 1 reputation and stack 2 Slime(S) to center deck"
}
func (m *SlimeMid) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 2)
	return cost
}
func (m *SlimeMid) OnPunish() {
	m.state.TakeDamage(2)
}
func (m *SlimeMid) OnSlain() {
	// m.state.AddResource(RESOURCE_NAME_EXPLORATION, 1)
	m.state.AddResource(RESOURCE_NAME_REPUTATION, 2)
	h := []Card{}
	for i := 0; i < 2; i++ {
		midSlime := NewSlimeSmall(m.state)
		h = append(h, &midSlime)
	}
	m.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, false, h...)
}
