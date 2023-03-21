package cards

type LichMageMonster struct {
	BaseMonster
	state       AbstractGamestate
	turnCounter int
}

func NewLichMageMonster(state AbstractGamestate) LichMageMonster {
	return LichMageMonster{state: state}
}

func (m *LichMageMonster) GetName() string {
	return "Lich Mage"
}
func (m *LichMageMonster) GetDescription() string {
	return "Add 1 Stun each turn, and additional 2 dmg every 3 turns, on slain unlocks Tomb of forgotten monarch: central chamber"
}
func (m *LichMageMonster) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 5)
	return cost
}
func (m *LichMageMonster) OnPunish() {
	m.turnCounter += 1
	curse1 := NewStunCurse(m.state)
	m.state.DiscardCard(&curse1, DISCARD_SOURCE_NAN)
	if m.turnCounter%3 == 0 {
		m.state.TakeDamage(2)
	}
}
