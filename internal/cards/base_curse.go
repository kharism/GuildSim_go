package cards

type BaseCurse struct{}

func (h *BaseCurse) GetName() string {
	return ""
}
func (h *BaseCurse) GetDescription() string {
	return ""
}
func (h *BaseCurse) GetCost() Cost {
	cost := NewCost()
	return cost
}

// return Hero,Area,Monster,Event etc
func (h *BaseCurse) GetCardType() CardType {
	return Curse
}
func (h *BaseCurse) OnAddedToHand() {}

// when played from hand, do this
func (h *BaseCurse) OnPlay() {}

// when explored, do this
func (h *BaseCurse) OnExplored() {}

// when slain, do this
func (h *BaseCurse) OnSlain() {}

// when discarded to cooldown pile, do this
func (h *BaseCurse) OnDiscarded() {}

func (h *BaseCurse) OnPunish() {

}
func (h *BaseCurse) Dispose(source string) {

}
