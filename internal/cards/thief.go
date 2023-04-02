package cards

type Thief struct {
	BaseHero
	gamestate AbstractGamestate
}

func (r *Thief) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}

func NewThief(gamestate AbstractGamestate) Thief {
	return Thief{gamestate: gamestate}
}

func (r *Thief) GetName() string {
	return "Thief"
}
func (r *Thief) GetDescription() string {
	return "gain 2 Exploration or peek top of center deck, if it is has a trap, disarm it"
}
func (r *Thief) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return cost
}
func (r *Thief) OnPlay() {
	if r.gamestate.GetBoolPicker().BoolPick("Gain 2 exploration?") {
		r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	} else {
		//r.gamestate.AddResource(RESOURCE_NAME_BLOCK, 5)
		topCenter := r.gamestate.PeekCenterCard()
		r.gamestate.GetDetailViewer().ShowDetail(topCenter)
		if _, ok := topCenter.(Trapper); ok {
			j := topCenter.(Trapper)
			j.Disarm()
		}

	}

}
