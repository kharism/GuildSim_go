package cards

type Cleric struct {
	BaseHero
	gamestate AbstractGamestate
}

func (r *Cleric) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func NewCleric(gamestate AbstractGamestate) Cleric {
	return Cleric{gamestate: gamestate}
}

func (r *Cleric) GetName() string {
	return "Cleric"
}
func (r *Cleric) GetDescription() string {
	return "Banish 1 non-monster from center row then draw 1 card"
}
func (r *Cleric) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return cost
}
func (r *Cleric) OnPlay() {
	centerCards := r.gamestate.GetCenterCard()
	availCard := []Card{}
	for _, i := range centerCards {
		_, isUnbanishable := i.(Unbanishable)
		if i.GetCardType() != Monster && !isUnbanishable {
			availCard = append(availCard, i)
		}
	}
	cardPicker := r.gamestate.GetCardPicker()
	banishIdx := cardPicker.PickCard(availCard, "Choose Card to banish")
	banishedCard := availCard[banishIdx]
	r.gamestate.RemoveCardFromCenterRow(banishedCard)
	r.gamestate.UpdateCenterCard(banishedCard)
	r.gamestate.BanishCard(banishedCard, DISCARD_SOURCE_CENTER)

	r.gamestate.Draw()
}
