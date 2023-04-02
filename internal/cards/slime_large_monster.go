package cards

type SlimeLarge struct {
	BaseMonster
	state       AbstractGamestate
	turnCounter int
}

func NewSlimeLarge(state AbstractGamestate) SlimeLarge {
	k := SlimeLarge{state: state}
	return k
}
func (m *SlimeLarge) GetName() string {
	return "Slime(L)"
}
func (m *SlimeLarge) GetDescription() string {
	return "deals 3 damage, Also cooldown 1 stun every 2 turns.on defeat, gains 2 reputation and stack 3 Slime(M) to center deck"
}
func (m *SlimeLarge) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 3)
	return cost
}
func (m *SlimeLarge) OnPunish() {
	m.state.TakeDamage(3)
	if m.turnCounter%3 == 0 {
		// m.state.TakeDamage(2)
		curse1 := NewStunCurse(m.state)
		m.state.DiscardCard(&curse1, DISCARD_SOURCE_NAN)
	}
}
func (m *SlimeLarge) OnSlain() {
	// m.state.AddResource(RESOURCE_NAME_EXPLORATION, 1)
	m.state.AddResource(RESOURCE_NAME_REPUTATION, 2)
	h := []Card{}
	for i := 0; i < 3; i++ {
		midSlime := NewSlimeMid(m.state)
		h = append(h, &midSlime)
	}
	m.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, false, h...)
}
