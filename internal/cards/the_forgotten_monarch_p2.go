package cards

type ForgottenMonarchP2 struct {
	BaseMonster
	state       AbstractGamestate
	turnCounter int
	isDisarmed  bool
}

func NewForgottenMonarchP2(state AbstractGamestate) ForgottenMonarchP2 {
	return ForgottenMonarchP2{state: state}
}

func (m *ForgottenMonarchP2) GetName() string {
	return "ForgottenMonarchP2"
}
func (m *ForgottenMonarchP2) GetDescription() string {
	return "Trap: Stack 1 Stun on main deck. On Punish: take 9 damage at end of turn.  Reward: gain 10 reputation"
}
func (m *ForgottenMonarchP2) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 12)
	return cost
}
func (m *ForgottenMonarchP2) Trap() {
	s := NewStunCurse(m.state)
	m.state.StackCards(DISCARD_SOURCE_NAN, &s)
}
func (m *ForgottenMonarchP2) Unbanishable() {

}
func (m *ForgottenMonarchP2) Unshuffleable() {

}
func (m *ForgottenMonarchP2) IsDisarmed() bool {
	return m.isDisarmed
}
func (m *ForgottenMonarchP2) Disarm() {
	m.isDisarmed = true
}
func (m *ForgottenMonarchP2) OnDisarm() {

}
func (m *ForgottenMonarchP2) OnSlain() {
	m.state.AddResource(RESOURCE_NAME_REPUTATION, 10)
	data := map[string]interface{}{}
	data[EVENT_ATTR_BOSS_DEFEATED_COUNT] = 1
	m.state.RemoveQuest("Explore Tomb of Forgotten monarch and defeat whatever lies there")
	m.state.NotifyListener(EVENT_BOSS_DEFEATED, data)
}
func (m *ForgottenMonarchP2) OnPunish() {
	m.turnCounter += 1
	m.state.TakeDamage(9)
	if m.turnCounter%2 == 0 {
		curse1 := NewStunCurse(m.state)
		m.state.DiscardCard(&curse1, DISCARD_SOURCE_NAN)
	}
}
