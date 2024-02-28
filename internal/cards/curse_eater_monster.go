package cards

type CurseEaterMonster struct {
	BaseMonster
	gamestate AbstractGamestate
}
type CurseEaterDagger struct {
	BaseItem
	gamestate AbstractGamestate
}

func (c *CurseEaterDagger) GetName() string {
	return "CurseEaterDagger"
}
func (c *CurseEaterDagger) GetDescription() string {
	return "If you play a curse card, gain 2 combat, 2 block, and draw 1"
}
func (c *CurseEaterDagger) OnAcquire() {
	action1 := NewAddResourceAction(c.gamestate, RESOURCE_NAME_COMBAT, 2)
	action2 := NewAddResourceAction(c.gamestate, RESOURCE_NAME_BLOCK, 2)
	action3 := NewDrawAction(c.gamestate)
	compAction := NewCompositeAction(c.gamestate, action1, action2, action3)
	filter := &CardFilter{Key: FILTER_TYPE, Op: Eq, Value: Curse}
	listener := NewCardPlayedListener(filter, compAction)
	c.gamestate.AttachListener(EVENT_CARD_PLAYED, listener)
}
func NewCurseEaterMonster(state AbstractGamestate) CurseEaterMonster {
	return CurseEaterMonster{gamestate: state}
}
func (r *CurseEaterMonster) GetName() string {
	return "CurseEaterMonster"
}
func (r *CurseEaterMonster) GetDescription() string {
	return "on punish: 3 damage. onslain: gain curse eater dagger"
}
func (r *CurseEaterMonster) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 4)
	return cost
}
func (m *CurseEaterMonster) OnPunish() {
	m.gamestate.TakeDamage(3)
}
func (r *CurseEaterMonster) OnSlain() {
	//r.gamestate.AddResource(RESOURCE_NAME_REPUTATION, 3)
	//r.gamestate.Draw()
	dagger := CurseEaterDagger{gamestate: r.gamestate}
	r.gamestate.AddItem(&dagger)
}
