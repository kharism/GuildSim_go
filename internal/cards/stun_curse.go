package cards

type StunCurse struct {
	BaseCurse
	state AbstractGamestate
}

func NewStunCurse(state AbstractGamestate) StunCurse {
	return StunCurse{state: state}
}
func (h *StunCurse) GetName() string {
	return "StunCurse"
}
func (h *StunCurse) GetDescription() string {
	return "this card do nothing on play"
}
func (h *StunCurse) GetCost() Cost {
	cost := NewCost()
	return cost
}
func (h *StunCurse) Dispose(source string) {
	h.state.DiscardCard(h, source)
}
