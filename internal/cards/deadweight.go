package cards

type Deadweight struct {
	BaseHero
	gamestate AbstractGamestate
}

func (r *Deadweight) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func NewDeadweight(gamestate AbstractGamestate) Deadweight {
	return Deadweight{gamestate: gamestate}
}

func (r *Deadweight) GetName() string {
	return "Deadweight"
}
func (r *Deadweight) GetDescription() string {
	return "if discarded, you either gain 3 combat or 3 exploration"
}
func (r *Deadweight) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}
func (r *Deadweight) OnDiscarded() {
	boolPicker := r.gamestate.GetBoolPicker()
	result := boolPicker.BoolPick("Gain 3 combat?")
	if result {
		r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 3)
	} else {
		r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	}
}
