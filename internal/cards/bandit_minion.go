package cards

type BanditMinion struct {
	BaseMonster
	state AbstractGamestate
}

func NewBanditMinion(state AbstractGamestate) BanditMinion {
	return BanditMinion{state: state}
}
func (b *BanditMinion) GetName() string {
	return "BanditMinion"
}
func (b *BanditMinion) GetDescription() string {
	return "On punish : 4 damage on end of turn."
}
func (b *BanditMinion) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}

func (b *BanditMinion) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 3)
	return cost
}
func (b *BanditMinion) OnSlain() {
	//b.state.TakeDamage(-2)
}
func (b *BanditMinion) OnPunish() {
	b.state.TakeDamage(4)
}
