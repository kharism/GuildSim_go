package cards

type Thief struct {
	BaseHero
	gamestate AbstractGamestate
}

func (r *Thief) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}

func NewThief(gamestate AbstractGamestate) Thief {
	return Thief{gamestate: gamestate}
}

func (r *Thief) GetName() string {
	return "Thief"
}
func (r *Thief) GetDescription() string {
	return "gain 2 Exploration or disarm trap in center row"
}
func (r *Thief) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return cost
}
func (r *Thief) OnPlay() {
	centerCard := r.gamestate.GetCenterCard()
	trapInCenter := []Card{}
	for _, c := range centerCard {
		if c.GetCardType() == Trap {
			trapInCenter = append(trapInCenter, c)
		}
	}
	if len(trapInCenter) == 0 {
		r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	} else {
		if r.gamestate.GetBoolPicker().BoolPick("Gain 2 exploration?") {
			r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 2)
		} else {
			//r.gamestate.AddResource(RESOURCE_NAME_BLOCK, 5)
			selectedIdx := r.gamestate.GetCardPicker().PickCard(trapInCenter, "Pick a card to disarm")
			selectedCard := trapInCenter[selectedIdx]
			// we do not disarm properly, since we do it for free.
			// r.gamestate.Disarm(selectedCard)
			selectedCard.(Trapper).OnDisarm()
			r.gamestate.RemoveCardFromCenterRow(selectedCard)
			r.gamestate.UpdateCenterCard(selectedCard)
			selectedCard.Dispose(DISCARD_SOURCE_CENTER)
			trapRemovedEvent := map[string]interface{}{EVENT_ATTR_TRAP_REMOVED: selectedCard}
			r.gamestate.NotifyListener(EVENT_TRAP_REMOVED, trapRemovedEvent)
			// d.BanishCard(c, cards.DISCARD_SOURCE_CENTER)
			// r.gamestate.Dispose(cards.DISCARD_SOURCE_CENTER)
		}
	}

}
