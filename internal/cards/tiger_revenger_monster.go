package cards

import "github/kharism/GuildSim_go/internal/observer"

type TigerRevenger struct {
	BaseMonster
	isDisarmed     bool
	state          AbstractGamestate
	defeatListener observer.Listener
}

func NewTigerRevenger(state AbstractGamestate) TigerRevenger {
	cardFilter := CardFilter{Key: FILTER_TYPE, Op: Eq, Value: Monster}
	actionLooseCombat := NewAddResourceAction(state, RESOURCE_NAME_COMBAT, -1)
	monsterDefeatedlistener := NewCardDefeatedListener(&cardFilter, actionLooseCombat)
	return TigerRevenger{state: state, defeatListener: monsterDefeatedlistener}
}
func (b *TigerRevenger) GetName() string {
	return "TigerRevenger"
}
func (b *TigerRevenger) GetDescription() string {
	return "Trap : Each time another monster is defeated loose 1 combat. reward : 5 reputation. Unshuffleable. Can't be disarmed"
}
func (b *TigerRevenger) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}
func (b *TigerRevenger) IsDisarmed() bool {
	return b.isDisarmed
}
func (b *TigerRevenger) OnDisarm() {

}
func (b *TigerRevenger) Disarm() {

}
func (b *TigerRevenger) Trap() {
	b.state.AttachListener(EVENT_CARD_DEFEATED, b.defeatListener)
}
func (b *TigerRevenger) Unshuffleable() {

}
func (b *TigerRevenger) Unbanishable() {

}
func (b *TigerRevenger) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 5)
	return cost
}
func (b *TigerRevenger) OnSlain() {
	b.state.RemoveListener(EVENT_CARD_DEFEATED, b.defeatListener)
	b.state.AddResource(RESOURCE_NAME_REPUTATION, 5)
	jj := NewDragonLord(b.state)
	b.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, true, &jj)
}
func (b *TigerRevenger) OnPunish() {
	b.state.TakeDamage(4)
}
