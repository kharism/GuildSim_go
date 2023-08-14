package cards

type FireMage struct {
	BaseHero
	gamestate AbstractGamestate
}

func (r *FireMage) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func NewFireMage(gamestate AbstractGamestate) FireMage {
	return FireMage{gamestate: gamestate}
}
func (r *FireMage) GetName() string {
	return "Fire Mage"
}
func (r *FireMage) GetDescription() string {
	return "Add 1 Combat point and draw 1 card"
}
func (r *FireMage) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}
func (r *FireMage) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 1)
	r.gamestate.Draw()
}
