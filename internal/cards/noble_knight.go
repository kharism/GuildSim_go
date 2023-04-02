package cards

type NobleKnight struct {
	BaseHero
	gamestate AbstractGamestate
}

func (r *NobleKnight) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func NewNobleKnight(gamestate AbstractGamestate) NobleKnight {
	return NobleKnight{gamestate: gamestate}
}

func (r *NobleKnight) GetName() string {
	return "NobleKnight"
}
func (r *NobleKnight) GetDescription() string {
	return "Gain combat equal to half of your reputation"
}
func (r *NobleKnight) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 5)
	return cost
}
func (r *NobleKnight) OnPlay() {
	rep := r.gamestate.GetCurrentResource().Detail[RESOURCE_NAME_REPUTATION]
	r.gamestate.AddResource(RESOURCE_NAME_COMBAT, rep/2)
}
