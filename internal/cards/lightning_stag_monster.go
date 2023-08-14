package cards

type LightingStag struct {
	BaseMonster
	gamestate AbstractGamestate
}

func NewLightingStag(state AbstractGamestate) LightingStag {
	g := LightingStag{gamestate: state}
	return g
}

func (ed *LightingStag) GetName() string {
	return "LightningStag"
}
func (ed *LightingStag) Dispose(source string) {
	ed.gamestate.DiscardCard(ed, source)
}
func (ed *LightingStag) GetDescription() string {
	return "Recruitable. Reward: recover 5 HP and cooldown a stun curse. On Play: gains 4 exploration and stack a stun curse"
}
func (ed *LightingStag) GetCost() Cost {
	c := NewCost()
	c.AddResource(RESOURCE_NAME_COMBAT, 3)
	return c
}
func (g *LightingStag) OnPunish() {

}
func (m *LightingStag) OnRecruit() {

}
func (m *LightingStag) OnPlay() {
	m.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 4)
	stunCurse := NewStunCurse(m.gamestate)
	m.gamestate.StackCards(DISCARD_SOURCE_NAN, &stunCurse)
}
func (g *LightingStag) OnSlain() {
	g.gamestate.TakeDamage(-5)
	stunCurse := NewStunCurse(g.gamestate)
	g.gamestate.DiscardCard(&stunCurse, DISCARD_SOURCE_NAN)
	// g.gamestate.AddResource(RESOURCE_NAME_REPUTATION, 1)
}
