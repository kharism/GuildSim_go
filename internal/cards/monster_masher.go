package cards

type MonsterMasher struct {
	BaseHero
	gamestate AbstractGamestate
}

func (r *MonsterMasher) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func NewMonsterMasher(gamestate AbstractGamestate) MonsterMasher {
	return MonsterMasher{gamestate: gamestate}
}

func (r *MonsterMasher) GetName() string {
	return "MonsterMasher"
}
func (r *MonsterMasher) GetDescription() string {
	return "gain 3 combat, and additional 2 for each monster in center row"
}
func (r *MonsterMasher) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 5)
	return cost
}
func (r *MonsterMasher) OnPlay() {
	cardsInCenter := r.gamestate.GetCenterCard()
	monsterCount := 0
	for _, i := range cardsInCenter {
		if i.GetCardType() == Monster {
			monsterCount += 1
		}
	}
	r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 3+2*monsterCount)
}
