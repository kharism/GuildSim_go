package cards

type Deadweight struct {
	BaseHero
	isDisposed bool
	gamestate  AbstractGamestate
}

func (r *Deadweight) Dispose(source string) {
	r.isDisposed = true
	r.gamestate.DiscardCard(r, source)
}
func NewDeadweight(gamestate AbstractGamestate) Deadweight {
	return Deadweight{gamestate: gamestate}
}
func (r *Deadweight) OnAddedToHand() {
	r.isDisposed = false
}
func (r *Deadweight) GetName() string {
	return "Deadweight"
}
func (r *Deadweight) GetDescription() string {
	return "if discarded, you either gain 3 combat or 3 exploration. On play: gain 5 block"
}
func (r *Deadweight) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}
func (r *Deadweight) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_BLOCK, 5)
}
func (r *Deadweight) OnDiscarded() {
	if !r.isDisposed {
		boolPicker := r.gamestate.GetBoolPicker()
		result := boolPicker.BoolPick("Gain 3 combat?")
		if result {
			r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 3)
		} else {
			r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 3)
		}
	}

}
