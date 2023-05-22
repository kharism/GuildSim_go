package cards

type WildBoar struct {
	BaseMonster
	state AbstractGamestate
}

func NewWildBoar(state AbstractGamestate) WildBoar {
	return WildBoar{state: state}
}
func (b *WildBoar) GetName() string {
	return "WildBoar"
}
func (b *WildBoar) GetDescription() string {
	return "On punish : 2 damage on end of turn. reward : recover 2 HP"
}
func (b *WildBoar) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}

func (b *WildBoar) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 2)
	return cost
}
func (b *WildBoar) OnSlain() {
	b.state.TakeDamage(-2)
}
func (b *WildBoar) OnPunish() {
	b.state.TakeDamage(2)
}
