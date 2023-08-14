package cards

type Direwolf struct {
	BaseMonster
	state AbstractGamestate
}

func (b *Direwolf) GetName() string {
	return "Direwolf"
}
func (b *Direwolf) GetDescription() string {
	return "on punish: 3 damage on end of turn. Reward: banish 1 card from your hand and if you do cooldown 1 MonsterMasher"
}

func (b *Direwolf) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}
func NewDirewolf(state AbstractGamestate) Direwolf {
	return Direwolf{state: state}
}
func (b *Direwolf) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 7)
	return cost
}
func (b *Direwolf) OnPunish() {
	b.state.TakeDamage(3)
}

// TODO
func (b *Direwolf) OnSlain() {
	cardInHand := b.state.GetCardInHand()
	if len(cardInHand) > 0 {
		idx := b.state.GetCardPicker().PickCard(cardInHand, "Pick a card to banish")
		selectedCard := cardInHand[idx]
		b.state.RemoveCardFromHandIdx(idx)
		b.state.BanishCard(selectedCard, DISCARD_SOURCE_HAND)
		monsterMasher := NewMonsterMasher(b.state)
		b.state.DiscardCard(&monsterMasher, DISCARD_SOURCE_NAN)
	}
	// b.state.AddResource(RESOURCE_NAME_REPUTATION, 2)
	// if b.state.GetBoolPicker().BoolPick("Draw a card?") {
	// 	b.state.Draw()
	// } else {
	// 	b.state.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	// }
}
