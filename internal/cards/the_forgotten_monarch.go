package cards

type ForgottenMonarch struct {
	BaseMonster
	state       AbstractGamestate
	turnCounter int
}

func NewForgottenMonarch(state AbstractGamestate) ForgottenMonarch {
	return ForgottenMonarch{state: state}
}

func (m *ForgottenMonarch) GetName() string {
	return "ForgottenMonarch"
}
func (m *ForgottenMonarch) GetDescription() string {
	return "On Punish: take 8 damage at end of turn. every 2 turn discard stun curse. Reward: gain 4 HP and stack ForgottenRagingMonarch"
}
func (m *ForgottenMonarch) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 8)
	return cost
}
func (m *ForgottenMonarch) Unbanishable() {

}
func (m *ForgottenMonarch) Unshuffleable() {

}
func (m *ForgottenMonarch) OnSlain() {
	m.state.TakeDamage(-2)
	phase2 := NewForgottenMonarchP2(m.state)
	m.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, false, &phase2)
}
func (m *ForgottenMonarch) OnPunish() {
	m.turnCounter += 1
	m.state.TakeDamage(8)
	if m.turnCounter%2 == 0 {
		curse1 := NewStunCurse(m.state)
		m.state.DiscardCard(&curse1, DISCARD_SOURCE_NAN)
	}
}
