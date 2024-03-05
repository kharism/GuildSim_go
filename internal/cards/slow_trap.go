package cards

type SlowTrap struct {
	BaseTrap
	state         AbstractGamestate
	limitResource *LimitResource
	isDisarmed    bool
}

func NewSlowTrap(state AbstractGamestate) SlowTrap {
	j := SlowTrap{state: state, limitResource: NewLimitResource(state, RESOURCE_NAME_EXPLORATION, 2, 1)}
	// j.limitDraw.AttachLimitDraw(state)
	return j
}
func (b *SlowTrap) GetName() string {
	return "SlowTrap"
}
func (b *SlowTrap) Disarm() {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}
func (b *SlowTrap) IsDisarmed() bool {
	return b.isDisarmed
}
func (b *SlowTrap) GetDescription() string {
	return "if you generate 2 or more exploration point, reduce it by 1.Unshuffleable.Unbanishable"
}
func (b *SlowTrap) Unshuffleable() {}
func (b *SlowTrap) Unbanishable()  {}
func (b *SlowTrap) Dispose(source string) {
	b.limitResource.DetachListener()
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}
func (b *SlowTrap) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 5)
	return cost
}
func (b *SlowTrap) OnDisarm() {
	relic := b.state.GenerateRandomPotion(RARITY_COMMON)
	b.state.AddItem(relic)
	// data := map[string]interface{}{}
	//b.state.NotifyListener(EVENT_MINIBOSS_DEFEATED, data)
}
