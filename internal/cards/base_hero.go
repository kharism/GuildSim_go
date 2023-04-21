package cards

type BaseHero struct{}

func (h *BaseHero) GetName() string {
	return ""
}
func (h *BaseHero) GetDescription() string {
	return ""
}
func (h *BaseHero) GetCost() Cost {
	cost := NewCost()
	return cost
}

// return Hero,Area,Monster,Event etc
func (h *BaseHero) GetCardType() CardType {
	return Hero
}

func (h *BaseHero) OnAddedToHand() {}

// when played from hand, do this
func (h *BaseHero) OnPlay() {}

// when explored, do this
func (h *BaseHero) OnExplored() {}

// when slain, do this
func (h *BaseHero) OnSlain()        {}
func (a *BaseHero) OnBanished()     {}
func (a *BaseHero) OnReturnToDeck() {}

// when discarded to cooldown pile, do this
func (h *BaseHero) OnDiscarded() {}

func (h *BaseHero) OnRecruit() {

}
func (h *BaseHero) OnAcquire() {

}
func (h *BaseHero) Dispose(source string) {

}

// all Hero or recruitable monster/things must implement this interface
// although it does nothing on recuit
type Recruitable interface {
	OnRecruit()
}
