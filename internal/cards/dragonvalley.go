package cards

type DragonValley struct {
	BaseArea
	state AbstractGamestate
}

func (d *DragonValley) GetName() string {
	return "DragonValley"
}
func (d *DragonValley) GetDescription() string {
	return "If you have dragonclaw you can explore this card for free. OnExplore: Shuffle Undead Dragon and other dragons"
}

func NewDragonValley(state AbstractGamestate) DragonValley {
	return DragonValley{state: state}
}
func (d *DragonValley) GetCost() Cost {
	hasDragonClaw := false
	dragonClaw := DragonClaw{}
	items := d.state.ListItems()
	for _, i := range items {
		if i.GetName() == dragonClaw.GetName() {
			hasDragonClaw = true
			break
		}
	}
	cost := NewCost()
	if hasDragonClaw {
		return cost
	} else {
		cost.AddResource(RESOURCE_NAME_EXPLORATION, 10)
		return cost
	}
}
func (d *DragonValley) OnExplored() {
	undeadDragon := NewUndeadDragon(d.state)
	allCard := []Card{&undeadDragon}

	for c := 0; c < 3; c++ {
		jj := NewVerdantDragon(d.state)
		allCard = append(allCard, &jj)
	}
	for c := 0; c < 3; c++ {
		jj := NewBackBurner(d.state)
		allCard = append(allCard, &jj)
	}
	for i := 0; i < 6; i++ {
		ll := NewBlackBomber(d.state)
		allCard = append(allCard, &ll)
	}
	d.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, true, allCard...)
}
