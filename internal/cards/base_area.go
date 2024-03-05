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
func (a *BaseArea) GetKeywords() []string {
	return []string{}
}

// return Hero,Area,Monster,Event etc
func (a *BaseArea) GetCardType() CardType {
	return Area
}
func (a *BaseArea) Dispose(source string) {

}
func (a *BaseArea) OnAcquire() {

}
func (a *BaseArea) OnBanished()     {}
func (a *BaseArea) OnReturnToDeck() {}

// when played from hand, do this
func (a *BaseArea) OnPlay() {}

// when explored, do this
func (a *BaseArea) OnExplored() {}

func (a *BaseArea) OnAddedToHand() {}

// when slain, do this
func (a *BaseArea) OnSlain() {}

// when discarded to cooldown pile, do this
func (a *BaseArea) OnDiscarded() {}

// a hack. any card implement this interface should not be legal to banish
// any card that list a card to Banish should skip the card if it implements this interface
type Unbanishable interface {
	Unbanishable()
}
