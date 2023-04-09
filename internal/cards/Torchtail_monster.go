package cards

type Torchtail struct {
	BaseMonster
	state AbstractGamestate
}

func NewTorchtail(state AbstractGamestate) Torchtail {
	k := Torchtail{state: state}
	return k
}
func (m *Torchtail) GetName() string {
	return "Torchtail"
}
func (m *Torchtail) GetDescription() string {
	return "deals 2 damage, Also cooldown 1 burn. On defeat: gains 2 reputation and recover 10 HP"
}
func (m *Torchtail) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 4)
	return cost
}
func (m *Torchtail) OnPunish() {
	m.state.TakeDamage(2)
	curse1 := NewDamageEndturnCurse(m.state)
	m.state.StackCards(DISCARD_SOURCE_NAN, &curse1)
}
func (m *Torchtail) OnSlain() {
	// m.state.AddResource(RESOURCE_NAME_EXPLORATION, 1)
	m.state.AddResource(RESOURCE_NAME_REPUTATION, 2)
	m.state.TakeDamage(-10)
	// m.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, 	false, h...)
}
