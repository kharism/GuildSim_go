package cards

type Scout struct {
	BaseHero
	state AbstractGamestate
}

func NewScout(state AbstractGamestate) Scout {
	return Scout{state: state}
}

func (r *Scout) Dispose(source string) {
	r.state.DiscardCard(r, source)
}

func (r *Scout) GetName() string {
	return "Scout"
}
func (r *Scout) GetDescription() string {
	return "Add 1 Exploration point and draw 1 card"
}
func (r *Scout) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}
func (r *Scout) OnPlay() {
	r.state.AddResource(RESOURCE_NAME_EXPLORATION, 1)
	r.state.Draw()
}
