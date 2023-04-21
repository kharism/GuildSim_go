package cards

import "fmt"

type SlimeRoom struct {
	BaseTrap
	state      AbstractGamestate
	limitDraw  *LimitDraw
	isDisarmed bool
}

func NewSlimeRoom(state AbstractGamestate) SlimeRoom {
	j := SlimeRoom{state: state, limitDraw: NewLimitDraw(state, 5)}
	// j.limitDraw.AttachLimitDraw(state)
	return j
}
func (b *SlimeRoom) GetName() string {
	return "SlimeRoom"
}
func (b *SlimeRoom) GetDescription() string {
	return "when enter center row: stack 1 slime(L) on center. continuous eff: you can't draw more than 5 cards each turn"
}
func (b *SlimeRoom) Dispose(source string) {
	b.limitDraw.DetachLimitDraw(b.state)
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}

func (b *SlimeRoom) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 4)
	return cost
}
func (b *SlimeRoom) OnDisarm() {
	relic := b.state.GenerateRandomRelic(RARITY_RARE)
	b.state.AddItem(relic)
	data := map[string]interface{}{}
	b.state.NotifyListener(EVENT_MINIBOSS_DEFEATED, data)
}

func (b *SlimeRoom) Trap() {
	if !b.isDisarmed {
		for i := 0; i < 1; i++ {
			newSlime := NewSlimeLarge(b.state)
			b.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, false, &newSlime)
		}
	}
	fmt.Println("AttachLimitDraw From slimeroom")
	b.limitDraw.AttachLimitDraw(b.state)
	b.state.AttachLegalCheck(ACTION_DRAW, b.limitDraw)
}
func (b *SlimeRoom) Unshuffleable() {}
func (b *SlimeRoom) IsDisarmed() bool {
	return b.isDisarmed
}
func (b *SlimeRoom) Disarm() {
	b.isDisarmed = true
	b.OnDisarm()
}
