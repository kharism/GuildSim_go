package cards

type BoulderTrap struct {
	BaseTrap
	state      AbstractGamestate
	isDisarmed bool
}

func NewBoulderTrap(state AbstractGamestate) BoulderTrap {
	return BoulderTrap{state: state}
}
func (b *BoulderTrap) GetName() string {
	return "BoulderTrap"
}
func (b *BoulderTrap) GetDescription() string {
	return "when enter center row: stack 1 card from your hand to main deck then stack 2 stun"
}
func (b *BoulderTrap) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}

func (b *BoulderTrap) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 4)
	return cost
}

func (b *BoulderTrap) Trap() {
	if !b.isDisarmed {
		cardinhand := b.state.GetCardInHand()
		idx := b.state.GetCardPicker().PickCard(cardinhand, "pick card to be stacked")
		stackedCard := cardinhand[idx]
		b.state.RemoveCardFromHandIdx(idx)
		b.state.StackCards(DISCARD_SOURCE_HAND, stackedCard)
		for i := 0; i < 2; i++ {
			newStun := NewStunCurse(b.state)
			b.state.StackCards(DISCARD_SOURCE_NAN, &newStun)
		}
	}
}
func (b *BoulderTrap) IsDisarmed() bool {
	return b.isDisarmed
}
func (b *BoulderTrap) Disarm() {
	b.isDisarmed = true
}
