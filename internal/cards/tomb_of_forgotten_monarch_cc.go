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
	return "Rewards: 1 Rare relic, release guardians, defeat these monsters to "
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

	thunderGolem := NewThunderGolem(a.state)
	rottingGolem := NewRottingGolem(a.state)

	cardFilter := &CardFilter{Key: FILTER_NAME, Op: In, Value: []string{thunderGolem.GetName(), rottingGolem.GetName()}}

	cardsAdded := []Card{&thunderGolem, &rottingGolem}
	pushCenterDeckAction := NewPushCenterDeckAction(a.state, cardsAdded, true)
	removeEventListenerAction := NewRemoveEventListenerAction(a.state, EVENT_CARD_DEFEATED, nil)
	compositeAction := NewCompositeAction(a.state, pushCenterDeckAction, removeEventListenerAction)
	countDownAction := NewCountDownAction(len(cardsAdded), 1, compositeAction)
	guardiansDefeatedListener := NewCardDefeatedListener(cardFilter, countDownAction)
	removeEventListenerAction.(*RemoveEventListenerAction).listener = guardiansDefeatedListener

	a.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, true, cardsAdded...)
}
