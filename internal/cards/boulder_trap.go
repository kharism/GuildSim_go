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
	return "when enter center row: stack 1 card from your hand to main deck then stack 1 stun"
}
func (b *BoulderTrap) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}
func (b *BoulderTrap) OnDisarm() {

}
func (b *BoulderTrap) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}

func (b *BoulderTrap) Trap() {
	if !b.isDisarmed {
		b.state.MutexLock()
		defer b.state.MutexUnlock()
		cardinhand := b.state.GetCardInHand()
		if len(cardinhand) > 0 {
			idx := b.state.GetCardPicker().PickCard(cardinhand, "pick card to be stacked")
			if idx >= 0 {
				stackedCard := cardinhand[idx]
				b.state.RemoveCardFromHandIdx(idx)
				b.state.StackCards(DISCARD_SOURCE_HAND, stackedCard)
			}

		}
		for i := 0; i < 1; i++ {
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
