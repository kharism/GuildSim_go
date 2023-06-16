package cards

type ShieldBasher struct {
	BaseHero
	state AbstractGamestate
}

func NewShieldBasher(state AbstractGamestate) ShieldBasher {
	j := ShieldBasher{state: state}
	return j
}

func (h *ShieldBasher) GetName() string {
	return "Shield basher"
}
func (h *ShieldBasher) GetDescription() string {
	return "gain combat equal to your block"
}
func (h *ShieldBasher) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return cost
}
func (h *ShieldBasher) Dispose(source string) {
	h.state.DiscardCard(h, source)
}

func (h *ShieldBasher) OnPlay() {
	// h.state.MutexLock()
	res := h.state.GetCurrentResource()
	blockCounter := res.Detail[RESOURCE_NAME_BLOCK]
	h.state.AddResource(RESOURCE_NAME_COMBAT, blockCounter)
	// h.state.MutexUnlock()
}
