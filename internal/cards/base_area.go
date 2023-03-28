package cards

type BaseArea struct{}

func (a *BaseArea) GetName() string {
	return ""
}
func (a *BaseArea) GetDescription() string {
	return ""
}
func (a *BaseArea) GetCost() Cost {
	return NewCost()
}

// return Hero,Area,Monster,Event etc
func (a *BaseArea) GetCardType() CardType {
	return Area
}
func (a *BaseArea) Dispose(source string) {

}
func (a *BaseArea) OnAcquire() {

}

// when played from hand, do this
func (a *BaseArea) OnPlay() {}

// when explored, do this
func (a *BaseArea) OnExplored() {}

func (a *BaseArea) OnAddedToHand() {}

// when slain, do this
func (a *BaseArea) OnSlain() {}

// when discarded to cooldown pile, do this
func (a *BaseArea) OnDiscarded() {}
