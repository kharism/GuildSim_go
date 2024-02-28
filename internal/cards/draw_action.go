package cards

type DrawAction struct {
	state AbstractGamestate
}

func (d *DrawAction) DoAction() {
	d.state.Draw()
}
func NewDrawAction(a AbstractGamestate) *DrawAction {
	return &DrawAction{state: a}
}
