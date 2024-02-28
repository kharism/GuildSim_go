package cards

import "fmt"

type BlackBomber struct {
	BaseMonster
	state    AbstractGamestate
	Overlays []Card
}

func NewBlackBomber(state AbstractGamestate) BlackBomber {
	return BlackBomber{state: state}
}
func (b *BlackBomber) GetName() string {
	return "BlackBomber"
}
func (b *BlackBomber) GetDescription() string {
	return "On punish : 2 damage on end of turn and overlay 1 goblin minion. reward : gain 1 smoke bomb"
}
func (b *BlackBomber) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}
func (b *BlackBomber) HasOverlayCard() bool {
	return len(b.Overlays) > 0
}
func (r *BlackBomber) AttachOverlayCard(newCard Card) {
	r.Overlays = append(r.Overlays, newCard)
	data := map[string]interface{}{}
	data[EVENT_ATTR_ADD_OVERLAY_BASE_CARD] = r
	data[EVENT_ATTR_ADD_OVERLAY_ADDED_CARD] = newCard
	r.state.NotifyListener(EVENT_ADD_OVERLAY, data)
}
func (b *BlackBomber) GetOverlay() []Card {
	return b.Overlays
}

// detach the top card of this stack
func (r *BlackBomber) Detach() {
	topCard := r.Overlays[len(r.Overlays)-1]
	r.Overlays = r.Overlays[:len(r.Overlays)-1]
	data := map[string]interface{}{}
	data[EVENT_ATTR_ADD_OVERLAY_BASE_CARD] = r
	data[EVENT_ATTR_ADD_OVERLAY_ADDED_CARD] = topCard
	r.state.NotifyListener(EVENT_REMOVE_OVERLAY, data)
}
func (r *BlackBomber) GetCost() Cost {
	cost := NewCost()
	// cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	if r.HasOverlayCard() {
		return r.Overlays[len(r.Overlays)-1].GetCost()
	} else {
		cost.AddResource(RESOURCE_NAME_COMBAT, 2)
	}

	return cost
}
func (b *BlackBomber) OnSlain() {
	smokeBomb := SmokeBomb{state: b.state}
	b.state.AddItem(&smokeBomb)
}
func (b *BlackBomber) OnPunish() {
	b.state.TakeDamage(2)
	if b.HasOverlayCard() {
		pp := b.Overlays[len(b.Overlays)-1]
		if ll, ok := pp.(Punisher); ok {
			ll.OnPunish()
		}
	}

	hh := NewGoblinMinionMonster(b.state)
	b.AttachOverlayCard(&hh)
}

type SmokeBomb struct {
	BaseItem
	state AbstractGamestate
}

func (s *SmokeBomb) GetName() string {
	return "SmokeBomb"
}
func (h *SmokeBomb) GetDescription() string {
	return "Shuffle all cards on center row then replace them with new cards"
}

func (h *SmokeBomb) OnConsume() {
	cardsShuffledBack := 0
	cc := h.state.GetCenterCard()
	for i := len(cc) - 1; i >= 0; i-- {
		hh := cc[i]
		if _, ok := hh.(Unshuffleable); ok {
			continue
		}
		cardsShuffledBack++
		h.state.RemoveCardFromCenterRowIdx(i)
		h.state.AddCardToCenterDeck(DISCARD_SOURCE_CENTER, true, hh)
	}
	//h.state.CardsInCenterDeck.Shuffle()
	for i := 0; i < cardsShuffledBack; i++ {
		f := h.state.ReplaceCenterCard()
		fmt.Println("Replace center card with", f.GetName())
		h.state.AppendCenterCard(f)
		//h.state.CenterCards = append(d.CenterCards, f)
	}
}
