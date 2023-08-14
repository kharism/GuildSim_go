package cards

type ThunderGolem struct {
	BaseMonster
	turnCounter int
	state       AbstractGamestate
}

func NewThunderGolem(a AbstractGamestate) ThunderGolem {
	return ThunderGolem{state: a}
}

type CrystalKey struct {
	BaseItem
}

func (h *CrystalKey) GetName() string {
	return "CrystalKey"
}
func (h *CrystalKey) GetDescription() string {
	return "Open path to forgotten monarch"
}
func (t *ThunderGolem) GetName() string {
	return "ThunderGolem"
}
func (m *ThunderGolem) GetDescription() string {
	return "deals 4 damage also cooldown 2 shock every 2 turns. on defeat, gains 3 reputation and 1 key item"
}
func (m *ThunderGolem) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 7)
	return cost
}
func (m *ThunderGolem) OnPunish() {
	m.state.TakeDamage(4)
	m.turnCounter++
	if m.turnCounter%2 == 0 {
		for i := 0; i < 2; i++ {
			curse1 := NewShockCurse(m.state)
			m.state.DiscardCard(&curse1, DISCARD_SOURCE_NAN)
		}
	}

}
func (m *ThunderGolem) Unshuffleable() {}
func (m *ThunderGolem) Unbanishable()  {}
func (m *ThunderGolem) OnSlain() {
	m.state.AddResource(RESOURCE_NAME_REPUTATION, 3)
	crystalKey := CrystalKey{}
	m.state.AddItem(&crystalKey)
}
