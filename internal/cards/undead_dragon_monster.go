package cards

type UndeadDragon struct {
	BaseMonster
	state    AbstractGamestate
	Overlays []Card
}

func NewUndeadDragon(state AbstractGamestate) UndeadDragon {
	return UndeadDragon{state: state}
}
func (b *UndeadDragon) GetName() string {
	return "UndeadDragon"
}
func (b *UndeadDragon) GetDescription() string {
	return "On punish : 4 damage on end of turn and overlay 1 skeleton minion. reward : Banish 1 cards you have played"
}
func (b *UndeadDragon) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}
func (b *UndeadDragon) HasOverlayCard() bool {
	return len(b.Overlays) > 0
}
func (r *UndeadDragon) AttachOverlayCard(newCard Card) {
	r.Overlays = append(r.Overlays, newCard)
	data := map[string]interface{}{}
	data[EVENT_ATTR_ADD_OVERLAY_BASE_CARD] = r
	data[EVENT_ATTR_ADD_OVERLAY_ADDED_CARD] = newCard
	r.state.NotifyListener(EVENT_ADD_OVERLAY, data)
}
func (b *UndeadDragon) GetOverlay() []Card {
	return b.Overlays
}

// detach the top card of this stack
func (r *UndeadDragon) Detach() {
	topCard := r.Overlays[len(r.Overlays)-1]
	r.Overlays = r.Overlays[:len(r.Overlays)-1]
	data := map[string]interface{}{}
	data[EVENT_ATTR_ADD_OVERLAY_BASE_CARD] = r
	data[EVENT_ATTR_ADD_OVERLAY_ADDED_CARD] = topCard
	r.state.NotifyListener(EVENT_REMOVE_OVERLAY, data)
}
func (r *UndeadDragon) GetCost() Cost {
	cost := NewCost()
	// cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	if r.HasOverlayCard() {
		return r.Overlays[len(r.Overlays)-1].GetCost()
	} else {
		cost.AddResource(RESOURCE_NAME_COMBAT, 6)
	}

	return cost
}
func (b *UndeadDragon) OnSlain() {
	playedList := b.state.GetPlayedCards()
	if len(playedList) > 0 {
		selectedIdx := b.state.GetCardPicker().PickCard(playedList, "Pick card to banish")
		discardedCard := playedList[selectedIdx]
		b.state.RemoveCardFromHandIdx(selectedIdx)
		b.state.BanishCard(discardedCard, DISCARD_SOURCE_PLAYED)
	}
	data := map[string]interface{}{}
	data[EVENT_ATTR_CARD_DEFEATED] = b
	b.state.NotifyListener(EVENT_MINIBOSS_DEFEATED, data)

}
func (b *UndeadDragon) OnPunish() {
	b.state.TakeDamage(4)
	if b.HasOverlayCard() {
		pp := b.Overlays[len(b.Overlays)-1]
		if ll, ok := pp.(Punisher); ok {
			ll.OnPunish()
		}
	}

	hh := NewSkeletonMinion(b.state)
	b.AttachOverlayCard(&hh)
}
