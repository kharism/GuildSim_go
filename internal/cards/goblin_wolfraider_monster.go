package cards

type GoblinRaiderMonster struct {
	BaseMonster
	state AbstractGamestate
}

func NewGoblinWolfRaiderMonster(state AbstractGamestate) GoblinRaiderMonster {
	return GoblinRaiderMonster{state: state}

}
func (m *GoblinRaiderMonster) GetName() string {
	return "GoblinRaiderMonster"
}
func (m *GoblinRaiderMonster) GetDescription() string {
	return "2 damage per turn. Reward: banish 1 card on your cooldownpile"
}
func (m *GoblinRaiderMonster) GetCost() Cost {
	cost := NewCost()
	cost.Detail[RESOURCE_NAME_COMBAT] = 3
	return cost
}
func (m *GoblinRaiderMonster) OnPunish() {
	m.state.TakeDamage(2)
}
func (m *GoblinRaiderMonster) OnSlain() {
	cooldown := m.state.GetCooldownCard()
	if len(cooldown) > 0 {
		removeIdx := m.state.GetCardPicker().PickCard(cooldown, "Banish a card")
		removedCard := cooldown[removeIdx]
		m.state.RemoveCardFromCooldownIdx(removeIdx)
		m.state.BanishCard(removedCard, DISCARD_SOURCE_COOLDOWN)
	}
}
func (m *GoblinRaiderMonster) Dispose(source string) {
	m.state.DiscardCard(m, source)
}
