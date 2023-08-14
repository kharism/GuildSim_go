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
	return "Defeat 1 monster with cost less than 5 combat. If there aren't any, draw 2 cards"
}
func (r *MonsterMasher) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 5)
	return cost
}
func (r *MonsterMasher) OnPlay() {
	cardsInCenter := r.gamestate.GetCenterCard()
	legalTarget := []Card{}
	threshold := NewResource()
	threshold.AddResource(RESOURCE_NAME_COMBAT, 5)
	for _, i := range cardsInCenter {
		if i.GetCardType() == Monster {
			//monsterCount += 1
			cost := i.GetCost()
			if cost.IsEnough(threshold) {
				legalTarget = append(legalTarget, i)
			}
		}
	}
	if len(legalTarget) == 0 {
		r.gamestate.Draw()
		r.gamestate.Draw()
	} else {
		selected := r.gamestate.GetCardPicker().PickCard(legalTarget, "PIck one to defeat")
		selectedCard := legalTarget[selected]
		r.gamestate.BanishCard(selectedCard, DISCARD_SOURCE_CENTER)
		r.gamestate.RemoveCardFromCenterRow(selectedCard)
		selectedCard.OnSlain()

	}
	//r.gamestate.AddResource(RESOURCE_NAME_COMBAT, 3+2*monsterCount)
}
