package cards

import (
	"sort"
	"strings"
)

type AvalanceDragon struct {
	BaseMonster
	state AbstractGamestate
}

func NewAvalanceDragon(state AbstractGamestate) AvalanceDragon {
	return AvalanceDragon{state: state}
}
func (b *AvalanceDragon) GetName() string {
	return "AvalanceDragon"
}
func (b *AvalanceDragon) GetDescription() string {
	return "on punish:take 3 damage and discard 2 frozen curse. Reward: shuffle all cards in cooldown pile and draw 2 cards"
}
func (b *AvalanceDragon) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 6)
	return cost
}
func (b *AvalanceDragon) OnPunish() {
	b.state.TakeDamage(3)
	for i := 0; i < 2; i++ {
		cc := NewFreezeCurse(b.state)
		b.state.DiscardCard(&cc, DISCARD_SOURCE_NAN)
	}
}
func (b *AvalanceDragon) OnSlain() {
	discarded := b.state.GetCooldownCard()
	for i := len(discarded) - 1; i >= 0; i-- {
		d := discarded[i]
		b.state.RemoveCardFromCooldown(d)
		b.state.GetMainDeck().Stack(d)
	}
	//sort by name to speed up stacking process
	sort.Slice(discarded, func(i, j int) bool {
		iName := discarded[i].GetName()
		jName := discarded[j].GetName()
		return strings.Compare(iName, jName) < 0
	})
	//(DISCARD_SOURCE_COOLDOWN, discarded...)
	b.state.ShuffleMainDeck()
	b.state.Draw()
	b.state.Draw()
}
