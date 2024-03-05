package cards

type PotionVendor struct {
	BaseHero
	state      AbstractGamestate
	isPlayed   bool
	isDisarmed bool
	Overlays   []Card
}

func (r *PotionVendor) Dispose(source string) {
	if r.isPlayed {
		r.state.BanishCard(r, DISCARD_SOURCE_PLAYED)
	} else {
		r.state.DiscardCard(r, source)
	}
}

func NewPotionVendor(gamestate AbstractGamestate) PotionVendor {
	return PotionVendor{state: gamestate, Overlays: []Card{}}
}

func (r *PotionVendor) GetName() string {
	return "PotionVendor"
}
func (r *PotionVendor) GetDescription() string {
	return "on recruit: pick 1 potion out of 3 random potion. On play; generate 1 potion, then banish this card"
}
func (r *PotionVendor) GetCost() Cost {
	cost := NewCost()
	// cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	if r.HasOverlayCard() {
		return r.Overlays[len(r.Overlays)-1].GetCost()
	}
	return cost
}
func (r *PotionVendor) HasOverlayCard() bool {
	return len(r.Overlays) > 0
}
func (r *PotionVendor) AttachOverlayCard(newCard Card) {
	r.Overlays = append(r.Overlays, newCard)
	data := map[string]interface{}{}
	data[EVENT_ATTR_ADD_OVERLAY_BASE_CARD] = r
	data[EVENT_ATTR_ADD_OVERLAY_ADDED_CARD] = newCard
	r.state.NotifyListener(EVENT_ADD_OVERLAY, data)
}
func (r *PotionVendor) Detach() {
	topCard := r.Overlays[len(r.Overlays)-1]
	r.Overlays = r.Overlays[:len(r.Overlays)-1]
	data := map[string]interface{}{}
	data[EVENT_ATTR_ADD_OVERLAY_BASE_CARD] = r
	data[EVENT_ATTR_ADD_OVERLAY_ADDED_CARD] = topCard
	r.state.NotifyListener(EVENT_REMOVE_OVERLAY, data)
}
func (r *PotionVendor) OnPunish() {
	if r.HasOverlayCard() {
		pp := r.Overlays[len(r.Overlays)-1]
		if ll, ok := pp.(Punisher); ok {
			ll.OnPunish()
		}
	}
}
func (b *PotionVendor) GetOverlay() []Card {
	return b.Overlays
}
func (r *PotionVendor) Trap() {
	for i := 0; i < 4-len(r.GetOverlay()); i++ {
		jj := NewBanditMinion(r.state)
		r.AttachOverlayCard(&jj)
	}
}
func (r *PotionVendor) IsDisarmed() bool {
	return r.isDisarmed
}
func (r *PotionVendor) Disarm() {
	r.isDisarmed = false
}
func (r *PotionVendor) OnDisarm() {
	r.OnRecruit()
}
func (r *PotionVendor) OnRecruit() {
	potions := []Card{}
	uniqueCheck := map[string]bool{}
	for i := 0; i < 3; {
		jj := r.state.GenerateRandomPotion(RARITY_COMMON | RARITY_RARE)
		if _, ok := uniqueCheck[jj.GetName()]; !ok {
			potions = append(potions, jj)
			i += 1
			uniqueCheck[jj.GetName()] = true
		}
	}

	selectedIdx := r.state.GetCardPicker().PickCard(potions, "Pick a potion")
	selectedPotion := potions[selectedIdx]
	r.state.AddItem(selectedPotion)
}
func (r *PotionVendor) OnPlay() {
	potions := []Card{}
	for i := 0; i < 3; i++ {
		jj := r.state.GenerateRandomPotion(RARITY_COMMON | RARITY_RARE)
		potions = append(potions, jj)
	}

	selectedIdx := r.state.GetCardPicker().PickCard(potions, "Pick a potion")
	selectedPotion := potions[selectedIdx]
	r.state.AddItem(selectedPotion)
}
