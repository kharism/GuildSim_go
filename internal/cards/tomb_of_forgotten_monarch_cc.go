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
	return "Rewards: 1 Rare relic"
}
func (a *TombMonarchCC) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 5)
	return cost
}
func (a *TombMonarchCC) OnExplored() {
	//a.state.AddResource(RESOURCE_NAME_MONEY, 100)
	relic := a.state.GenerateRandomRelic(RARITY_RARE)
	a.state.AddItem(relic)

	//
}
