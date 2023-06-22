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
	return "if you have 0 block gain 5 block, otherwise gain combat equal to your block"
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
	if blockCounter == 0 {
		h.state.AddResource(RESOURCE_NAME_BLOCK, 5)
	} else {
		h.state.AddResource(RESOURCE_NAME_COMBAT, blockCounter)
	}

	// h.state.MutexUnlock()
}
