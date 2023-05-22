package cards

type FreezeCurse struct {
	BaseCurse
	state AbstractGamestate
}

func NewFreezeCurse(state AbstractGamestate) FreezeCurse {
	return FreezeCurse{state: state}
}

func (d *FreezeCurse) GetName() string {
	return "FreezeCurse"
}
func (d *FreezeCurse) GetDescription() string {
	return "do nothing. Banish this on end turn"
}
func (d *FreezeCurse) Dispose(source string) {
	// d.state.TakeDamage(2)
	d.state.BanishCard(d, source)
}
