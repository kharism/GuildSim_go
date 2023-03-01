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

// when slain, do this
func (m *BaseMonster) OnSlain() {}

// when discarded to cooldown pile, do this
func (m *BaseMonster) OnDiscarded() {}

// all monster that do punishing move on end phase should implement this
type Punisher interface {
	OnPunish()
}
