package cards

type CarcassSlimeCurse struct {
	BaseCurse
	played bool
	state  AbstractGamestate
}

func NewCarcassSlimeCurse(state AbstractGamestate) CarcassSlimeCurse {
	return CarcassSlimeCurse{state: state}
}
func (d *CarcassSlimeCurse) GetName() string {
	return "CarcasSlimeCurse"
}
func (d *CarcassSlimeCurse) GetDescription() string {
	return "Punish: Take 8 damage. On discard: take 5 damage. On play: stack 1 card from your hand, then banish this card on end turn instead"
}
func (b *CarcassSlimeCurse) OnPlay() {
	cardinhand := b.state.GetCardInHand()
	if len(cardinhand) > 0 {
		idx := b.state.GetCardPicker().PickCard(cardinhand, "pick card to be stacked")
		if idx >= 0 {
			stackedCard := cardinhand[idx]
			b.state.RemoveCardFromHandIdx(idx)
			b.state.StackCards(DISCARD_SOURCE_HAND, stackedCard)
			b.played = true
		}
	}

}
func (d *CarcassSlimeCurse) OnPunish() {
	d.state.TakeDamage(8)
}
func (d *CarcassSlimeCurse) OnDiscarded() {
	d.state.TakeDamage(5)
}
func (d *CarcassSlimeCurse) Dispose(source string) {
	// d.state.TakeDamage(2)
	if d.played {
		d.state.BanishCard(d, source)
	} else {
		d.state.DiscardCard(d, source)
	}

}
func (d *CarcassSlimeCurse) OnAddedToHand() {
	d.played = false
}
