package cards

type MonsterSlayer struct {
	BaseHero
	gamestate AbstractGamestate
}

func NewMonsterSlayer(state AbstractGamestate) MonsterSlayer {
	return MonsterSlayer{gamestate: state}
}

func (r *MonsterSlayer) Dispose() {
	r.gamestate.DiscardCard(r)
}

func (r *MonsterSlayer) GetName() string {
	return "Monster Slayer"
}
func (r *MonsterSlayer) GetDescription() string {
	return "Add 3 Combat point"
}
func (r *MonsterSlayer) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 4)
	return cost
}
func (r *MonsterSlayer) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 3)
}
