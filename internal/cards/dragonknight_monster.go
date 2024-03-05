package cards

type DragonKnight struct {
	BaseMonster
	gamestate AbstractGamestate
}

type DragonClaw struct {
	BaseItem
	gamestate AbstractGamestate
}

func (h *DragonClaw) GetName() string {
	return "DragonClaw"
}

func (h *DragonClaw) GetDescription() string {
	return "Gains 2 combat at the start of your turn for each card in discard pile"
}
func (h *DragonKnight) GetKeywords() []string {
	return []string{
		"DragonClaw :\nGains 1 combat at the start\nof your turn for each card in\ndiscard pile",
	}
}
func (h *DragonClaw) OnAcquire() {
	action := NewAddResourceDynamicAction(h.gamestate, RESOURCE_NAME_COMBAT, func() int {
		return len(h.gamestate.GetCooldownCard())
	})
	listener := NewBasicAction(action)
	h.gamestate.AttachListener(EVENT_START_OF_TURN, listener)
}

func NewDragonKnight(state AbstractGamestate) DragonKnight {
	return DragonKnight{gamestate: state}
}

func (b *DragonKnight) GetName() string {
	return "DragonKnight"
}
func (b *DragonKnight) GetDescription() string {
	return "reward: gain dragon claw. Punish: 3 damage"
}
func (b *DragonKnight) Dispose(source string) {
	b.gamestate.DiscardCard(b, source)
}
func (b *DragonKnight) OnPunish() {
	b.gamestate.TakeDamage(3)
}
func (b *DragonKnight) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 4)
	return cost
}
func (r *DragonKnight) OnSlain() {
	dragonClaw := DragonClaw{gamestate: r.gamestate}
	r.gamestate.AddItem(&dragonClaw)
	// skeletonDragon := NewUndeadDragon(r.gamestate)
	// r.gamestate.AddCardToCenterDeck(DISCARD_SOURCE_NAN, true, &skeletonDragon)
	// r.gamestate.AddResource(RESOURCE_NAME_REPUTATION, 3)
}
