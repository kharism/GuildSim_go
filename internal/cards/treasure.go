package cards

type Treasure struct {
	BaseTrap
	isDisarmed bool
	Overlays   []Card
	gamestate  AbstractGamestate
}

func (r *Treasure) Dispose(source string) {
	r.gamestate.BanishCard(r, source)
}

func NewTreasure(gamestate AbstractGamestate) Treasure {
	return Treasure{gamestate: gamestate, isDisarmed: false, Overlays: []Card{}}
}

func (r *Treasure) GetName() string {
	return "Treasure"
}
func (r *Treasure) GetDescription() string {
	return "Trap: overlay 2 goblin minion. Reward: gain 1 random rare relic"
}
func (r *Treasure) Trap() {
	for i := 0; i < 2; i++ {
		jj := NewGoblinMinionMonster(r.gamestate)
		r.AttachOverlayCard(&jj)
	}
}
func (r *Treasure) HasOverlayCard() bool {
	return len(r.Overlays) > 0
}
func (r *Treasure) AttachOverlayCard(newCard Card) {
	r.Overlays = append(r.Overlays, newCard)
	data := map[string]interface{}{}
	data[EVENT_ATTR_ADD_OVERLAY_BASE_CARD] = r
	data[EVENT_ATTR_ADD_OVERLAY_ADDED_CARD] = newCard
	r.gamestate.NotifyListener(EVENT_ADD_OVERLAY, data)
}
func (r *Treasure) Detach() {
	topCard := r.Overlays[len(r.Overlays)-1]
	r.Overlays = r.Overlays[:len(r.Overlays)-1]
	data := map[string]interface{}{}
	data[EVENT_ATTR_ADD_OVERLAY_BASE_CARD] = r
	data[EVENT_ATTR_ADD_OVERLAY_ADDED_CARD] = topCard
	r.gamestate.NotifyListener(EVENT_REMOVE_OVERLAY, data)
}
func (r *Treasure) OnPunish() {
	if r.HasOverlayCard() {
		pp := r.Overlays[len(r.Overlays)-1]
		if ll, ok := pp.(Punisher); ok {
			ll.OnPunish()
		}
	}
}
func (b *Treasure) GetOverlay() []Card {
	return b.Overlays
}
func (b *Treasure) IsDisarmed() bool {
	return b.isDisarmed
}
func (b *Treasure) Disarm() {
	b.isDisarmed = true
}
func (r *Treasure) GetCost() Cost {
	cost := NewCost()
	// cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	if r.HasOverlayCard() {
		return r.Overlays[len(r.Overlays)-1].GetCost()
	}
	return cost

}
func (r *Treasure) OnDisarm() {
	newRelic := r.gamestate.GenerateRandomRelic(RARITY_RARE)
	r.gamestate.AddItem(newRelic)
}
