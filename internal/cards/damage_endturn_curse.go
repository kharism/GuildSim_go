package cards

type DamageEndturnCurse struct {
	BaseCurse
	state AbstractGamestate
}

func NewDamageEndturnCurse(state AbstractGamestate) DamageEndturnCurse {
	return DamageEndturnCurse{state: state}
}
func (d *DamageEndturnCurse) GetName() string {
	return "EndturnBurn"
}
func (d *DamageEndturnCurse) GetDescription() string {
	return "on end of turn, inflict 2 damage then banish this card"
}
func (d *DamageEndturnCurse) Dispose(source string) {
	d.state.TakeDamage(2)
	d.state.BanishCard(d, source)
}
