package cards

type BackBurner struct {
	BaseMonster
	state AbstractGamestate
}

func NewBackBurner(state AbstractGamestate) BackBurner {
	return BackBurner{state: state}
}

func (v *BackBurner) GetName() string {
	return "BackBurner"
}
func (v *BackBurner) GetDescription() string {
	return "on punish: 1 damage per card played this turn and discard 1 burn curse. Reward: Optional banish 1 card in discard"
}
func (v *BackBurner) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 5)
	return cost
}
func (v *BackBurner) OnPunish() {
	cardsPlayed := v.state.GetPlayedCards()
	v.state.TakeDamage(len(cardsPlayed))
	burnCurse := NewDamageEndturnCurse(v.state)
	v.state.DiscardCard(&burnCurse, DISCARD_SOURCE_NAN)
}
func (v *BackBurner) OnSlain() {
	cardPicker := v.state.GetCardPicker()
	cooldown := v.state.GetCooldownCard()
	if len(cooldown) > 0 {
		removeIdx := cardPicker.PickCardOptional(cooldown, "you can pick 1 to banish")
		if removeIdx > 0 {
			removedCard := cooldown[removeIdx]
			v.state.RemoveCardFromCooldownIdx(removeIdx)
			v.state.BanishCard(removedCard, DISCARD_SOURCE_COOLDOWN)
		}
	}
}
