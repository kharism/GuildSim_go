package cards

type SlimeRoom struct {
	BaseTrap
	state      AbstractGamestate
	limitDraw  *LimitDraw
	isDisarmed bool
}

func NewSlimeRoom(state AbstractGamestate) SlimeRoom {
	return SlimeRoom{state: state}
}
func (b *SlimeRoom) GetName() string {
	return "SlimeRoom"
}
func (b *SlimeRoom) GetDescription() string {
	return "when enter center row: stack 2 slime(L) on center. continuous eff: you can't draw more than 5 cards each turn"
}
func (b *SlimeRoom) Dispose(source string) {
	b.limitDraw.DetachLimitDraw(b.state)
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}

func (b *SlimeRoom) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 6)
	return cost
}

func (b *SlimeRoom) Trap() {
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
		b.limitDraw = NewLimitDraw(b.state, 5)
		b.limitDraw.AttachLimitDraw(b.state)
	}
}
func (b *SlimeRoom) IsDisarmed() bool {
	return b.isDisarmed
}
func (b *SlimeRoom) Disarm() {
	b.isDisarmed = true
}
