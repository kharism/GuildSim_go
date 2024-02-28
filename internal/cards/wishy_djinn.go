package cards

type WishyDjinn struct {
	BaseMonster
	gamestate AbstractGamestate
}

func (b *WishyDjinn) GetName() string {
	return "WishyDjinn"
}
func (b *WishyDjinn) GetDescription() string {
	return "recruitable. Take 2 damage and draw 3 cards"
}
func (b *WishyDjinn) Dispose(source string) {
	b.gamestate.DiscardCard(b, DISCARD_SOURCE_PLAYED)
}
func NewWishyDjinn(state AbstractGamestate) WishyDjinn {
	b := WishyDjinn{gamestate: state}
	return b
}
func (b *WishyDjinn) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 2)
	return cost
}
func (r *WishyDjinn) OnSlain() {
	r.gamestate.AddResource(RESOURCE_NAME_REPUTATION, 3)
}
func (r *WishyDjinn) OnRecruit() {

}
func (r *WishyDjinn) OnPlay() {
	r.gamestate.TakeDamage(2)
	r.gamestate.Draw()
	r.gamestate.Draw()
	r.gamestate.Draw()
}
