package cards

type TombMonarchCC struct {
	BaseArea
	state AbstractGamestate
}

func NewTombMonarchCC(state AbstractGamestate) TombMonarchCC {
	h := TombMonarchCC{state: state}
	return h
}
func (a *TombMonarchCC) GetName() string {
	return "TombMonarchCC"
}
func (a *TombMonarchCC) GetDescription() string {
	return "Rewards: 500 Money, Unlock strange labrynth and put them into center deck"
}
func (a *TombMonarchCC) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 5)
	return cost
}
func (a *TombMonarchCC) OnExplored() {
	a.state.AddResource(RESOURCE_NAME_MONEY, 100)
}
