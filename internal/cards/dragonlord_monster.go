package cards

type DragonLord struct {
	BaseMonster
	gamestate AbstractGamestate
	TurnCount int
	Overlays  []Card
}

func NewDragonLord(state AbstractGamestate) DragonLord {
	return DragonLord{gamestate: state, Overlays: []Card{}}
}
func (b *DragonLord) GetName() string {
	return "DragonLord"
}
func (b *DragonLord) GetDescription() string {
	return "Trap: overlay 3 DragonLordGuard. On punish: 5 Damage and stack stun card.Every 2 turns overlay DragonLordGuardian"
}
func (r *DragonLord) OnPunish() {
	r.gamestate.TakeDamage(5)
	stun := NewStunCurse(r.gamestate)
	r.TurnCount += 1
	if r.TurnCount%2 == 0 {
		hh := NewDragonLordGuard(r.gamestate)
		r.AttachOverlayCard(&hh)
	}
	r.gamestate.StackCards(DISCARD_SOURCE_NAN, &stun)
	if r.HasOverlayCard() {
		pp := r.Overlays[len(r.Overlays)-1]
		if ll, ok := pp.(Punisher); ok {
			ll.OnPunish()
		}
	}
}
func (b *DragonLord) Dispose(source string) {
	b.gamestate.BanishCard(b, DISCARD_SOURCE_CENTER)
}
func (b *DragonLord) HasOverlayCard() bool {
	return len(b.Overlays) > 0
}
func (b *DragonLord) GetKeywords() []string {
	return []string{"DragonLordGuard:\n5 combat;Deals 5 dmg"}
}
func (r *DragonLord) OnSlain() {
	data := map[string]interface{}{}
	data[EVENT_ATTR_BOSS_DEFEATED_COUNT] = 2
	r.gamestate.NotifyListener(EVENT_BOSS_DEFEATED, data)
}
func (r *DragonLord) AttachOverlayCard(newCard Card) {
	r.Overlays = append(r.Overlays, newCard)
	data := map[string]interface{}{}
	data[EVENT_ATTR_ADD_OVERLAY_BASE_CARD] = r
	data[EVENT_ATTR_ADD_OVERLAY_ADDED_CARD] = newCard
	r.gamestate.NotifyListener(EVENT_ADD_OVERLAY, data)
}
func (b *DragonLord) GetOverlay() []Card {
	return b.Overlays
}

// detach the top card of this stack
func (r *DragonLord) Detach() {
	topCard := r.Overlays[len(r.Overlays)-1]
	r.Overlays = r.Overlays[:len(r.Overlays)-1]
	data := map[string]interface{}{}
	data[EVENT_ATTR_ADD_OVERLAY_BASE_CARD] = r
	data[EVENT_ATTR_ADD_OVERLAY_ADDED_CARD] = topCard
	r.gamestate.NotifyListener(EVENT_REMOVE_OVERLAY, data)
}
func (r *DragonLord) Trap() {
	for i := 0; i < 4; i++ {
		jj := NewDragonLordGuard(r.gamestate)
		r.AttachOverlayCard(&jj)
	}
}
func (r *DragonLord) IsDisarmed() bool {
	return false
}
func (r *DragonLord) Disarm() {
	// r.isDisarmed = false
}
func (r *DragonLord) OnDisarm() {

}
func (r *DragonLord) GetCost() Cost {
	cost := NewCost()
	// cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	if r.HasOverlayCard() {
		return r.Overlays[len(r.Overlays)-1].GetCost()
	} else {
		cost.AddResource(RESOURCE_NAME_COMBAT, 10)
	}

	return cost
}

type DragonLordGuard struct {
	BaseMonster
	gamestate AbstractGamestate
}

func NewDragonLordGuard(state AbstractGamestate) DragonLordGuard {
	k := DragonLordGuard{gamestate: state}
	return k
}
func (m *DragonLordGuard) GetName() string {
	return "DragonLordGuard"
}
func (m *DragonLordGuard) GetDescription() string {
	return "Deals 5 dmg"
}
func (m *DragonLordGuard) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 5)
	return cost
}
func (m *DragonLordGuard) OnPunish() {
	m.gamestate.TakeDamage(5)
}
func (m *DragonLordGuard) OnSlain() {
}
