package cards

type DamageDrawCurse struct {
	BaseCurse
	state AbstractGamestate
}

func NewDamageDrawCurse(state AbstractGamestate) DamageDrawCurse {
	return DamageDrawCurse{state: state}
}
func (d *DamageDrawCurse) GetName() string {
	return "Damage Draw Curse"
}
func (d *DamageDrawCurse) GetDescription() string {
	return "If drawn: Inflict 2 damage then banish this card"
}

// when added to hand do this
func (d *DamageDrawCurse) OnAddedToHand() {
	d.state.TakeDamage(2)
	d.state.RemoveCardFromHand(d)
	d.state.BanishCard(d, DISCARD_SOURCE_HAND)
}
