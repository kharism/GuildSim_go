package cards

type RottingGolem struct {
	BaseMonster
	state AbstractGamestate
}

func NewRottingGolem(a AbstractGamestate) RottingGolem {
	return RottingGolem{state: a}
}
func (t *RottingGolem) GetName() string {
	return "RottingGolem"
}
func (m *RottingGolem) GetDescription() string {
	return "stack carcass slime curse,  on defeat, gains 1 key item"
}
func (m *RottingGolem) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 7)
	return cost
}
func (d *RottingGolem) OnPunish() {
	//d.state.TakeDamage(8)
	curse := NewCarcassSlimeCurse(d.state)
	d.state.StackCards(DISCARD_SOURCE_NAN, &curse)
}
