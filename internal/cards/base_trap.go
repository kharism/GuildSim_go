package cards

// Not sure whether I should put in trap class or make some cards that starts in center dec
// have 'onEnterCenterRow' method to mimic 'fate' mechanic in ascension. But screw it I'll
// make one just to see what would happened
type BaseTrap struct {
}

func (b *BaseTrap) GetName() string {
	return "BaseTrap"
}
func (b *BaseTrap) GetDescription() string {
	return ""
}

func (b *BaseTrap) GetCost() Cost {
	return NewCost()
}

// return Hero,Area,Monster,Event etc
func (b *BaseTrap) GetCardType() CardType {
	return Trap
}

// when played from hand, do this
func (b *BaseTrap) OnPlay() {}

// when explored, do this
func (b *BaseTrap) OnExplored() {}

// when slain, do this
func (b *BaseTrap) OnSlain() {}

// when discarded to cooldown pile, do this
func (b *BaseTrap) OnDiscarded() {}

// when added to hand do this
func (b *BaseTrap) OnAddedToHand() {}

// when a card is added to item list
func (b *BaseTrap) OnAcquire() {}

// get rid of this card, you either send this to discard pile or banished pile
func (b *BaseTrap) Dispose(source string) {}

func (b *BaseTrap) Trap() {}

// gamestate checks whether a card drawn implement this method
// if so, execute the Trap() method
// the IsDisarmed() and Disarm() is there so we have the option to
// disarm the trap before it harm us
type Trapper interface {
	Trap()
	IsDisarmed() bool
	Disarm()
}
