package cards

type DartTrap struct {
	BaseTrap
	state      AbstractGamestate
	isDisarmed bool
}

func NewDartTrap(state AbstractGamestate) DartTrap {
	return DartTrap{state: state}
}
func (b *DartTrap) GetName() string {
	return "DartTrap"
}
func (b *DartTrap) GetDescription() string {
	return "trap: take 1 damage for each card you played, punish: take 1 damage for each card you have played. Disarm: gain 1 common potion"
}
func (b *DartTrap) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}

func (b *DartTrap) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 6)
	return cost
}
func (b *DartTrap) OnPunish() {
	playedCards := b.state.GetPlayedCards()
	b.state.TakeDamage(len(playedCards))
}

func (b *DartTrap) Trap() {
	if !b.isDisarmed {
		playedCards := b.state.GetPlayedCards()
		b.state.TakeDamage(len(playedCards))
	}
}
func (b *DartTrap) IsDisarmed() bool {
	return b.isDisarmed
}
func (b *DartTrap) Disarm() {
	b.isDisarmed = true
}
func (b *DartTrap) OnDisarm() {
	relic := b.state.GenerateRandomPotion(RARITY_COMMON)
	b.state.AddItem(relic)
}
