package cards

type WolfShaman struct {
	BaseHero
	gamestate    AbstractGamestate
	curseRemoval *CurseRemovalListener
}

func NewWolfShaman(state AbstractGamestate) WolfShaman {
	this := WolfShaman{BaseHero: BaseHero{}, curseRemoval: &CurseRemovalListener{state: state}}
	this.gamestate = state
	return this
}
func (r *WolfShaman) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func (r *WolfShaman) GetName() string {
	return "WolfShaman"
}
func (r *WolfShaman) GetCost() Cost {
	j := NewCost()
	j.AddResource(RESOURCE_NAME_EXPLORATION, 4)
	return j
}
func (r *WolfShaman) GetDescription() string {
	return "the 1st curse stacked to deck is banished"
}

type CurseRemovalListener struct {
	state AbstractGamestate
}

func (t *CurseRemovalListener) DoAction(data map[string]interface{}) {
	cardPlayed := data[EVENT_ATTR_CARD_STACKED].(Card)
	//fmt.Println("Detect center draw")
	if cardPlayed.GetCardType() == Curse {
		// check top of deck
		if t.state.GetMainDeck().cards[0] == cardPlayed {
			newList := t.state.GetMainDeck().List()[1:]
			t.state.GetMainDeck().SetList(newList)
			// invoke
			data2 := map[string]interface{}{}
			data2[EVENT_ATTR_CARD_BANISHED] = cardPlayed
			data2[EVENT_ATTR_DISCARD_SOURCE] = DISCARD_SOURCE_MAIN_DECK
			t.state.NotifyListener(EVENT_CARD_BANISHED, data2)
			t.state.RemoveListener(EVENT_CARD_STACKED, t)
		}
	}
}
func (r *WolfShaman) OnPlay() {
	removeEventListenerAction := NewRemoveEventListenerAction(r.gamestate, EVENT_END_OF_TURN, r.curseRemoval)
	removeEventListenerListener := NewBasicAction(removeEventListenerAction)
	removeEventListenerAction.(*RemoveEventListenerAction).AddListener(removeEventListenerListener)
	r.gamestate.AttachListener(EVENT_ATTR_CARD_STACKED, r.curseRemoval)
	r.gamestate.AttachListener(EVENT_END_OF_TURN, removeEventListenerListener)
}
