package cards

type BaseMonster struct{}

func (m *BaseMonster) GetName() string {
	return ""
}
func (m *BaseMonster) GetDescription() string {
	return ""
}
func (m *BaseMonster) GetCost() Cost {
	cost := NewCost()
	return cost
}

// return Hero,Area,Monster,Event etc
func (m *BaseMonster) GetCardType() CardType {
	return Monster
}

// when played from hand, do this
func (m *BaseMonster) OnPlay() {}

// when explored, do this
func (m *BaseMonster) OnExplored() {}

func (m *BaseMonster) OnAddedToHand() {}

func (a *BaseMonster) OnBanished()     {}
func (a *BaseMonster) OnReturnToDeck() {}

// when slain, do this
func (m *BaseMonster) OnSlain() {}

func (m *BaseMonster) Dispose(source string) {}

func (m *BaseMonster) OnAcquire() {}

// when discarded to cooldown pile, do this
func (m *BaseMonster) OnDiscarded() {}

// all monster that do punishing move on end phase should implement this
type Punisher interface {
	OnPunish()
}

func (a *BaseMonster) GetKeywords() []string {
	return []string{}
}

// a hack. any card implement this interface will not be shuffled back on end of turn
// when we can't defeat/explore/recruit cards on center row. They also cannot be returned
// by card eff such as winged lion
type Unshuffleable interface {
	Unshuffleable()
}
