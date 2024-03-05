package cards

import "golang.org/x/exp/rand"

type Firelake struct {
	BaseArea
	gamestate AbstractGamestate
}

func NewFirelake(state AbstractGamestate) Firelake {
	h := Firelake{gamestate: state}
	return h
}

func (ed *Firelake) GetName() string {
	return "Firelake"
}
func (ed *Firelake) GetDescription() string {
	return "Reward: Add 1 rare potion, unleash fiery creature into center deck. On punish: if you have 3 or less combat, take 2 damage"
}
func (ed *Firelake) GetCost() Cost {
	c := NewCost()
	c.AddResource(RESOURCE_NAME_EXPLORATION, 5)
	return c
}

// when played on hand, to this
func (ed *Firelake) OnPlay() {}
func (ed *Firelake) OnExplored() {
	// ed.gamestate.AddResource(RESOURCE_NAME_MONEY, 100)
	newPotion := ed.gamestate.GenerateRandomPotion(RARITY_RARE)
	ed.gamestate.AddItem(newPotion)
	// ed.gamestate.AddResource(RESOURCE_NAME_REPUTATION, 2)
	monsters := []Card{}
	for i := 0; i < 2; i++ {
		ll := NewAggroDjinn(ed.gamestate)
		monsters = append(monsters, &ll)
	}
	for i := 0; i < 4; i++ {
		ll := NewInfernalJester(ed.gamestate)
		monsters = append(monsters, &ll)
	}
	for i := 0; i < 4; i++ {
		ll := NewEmberEater(ed.gamestate)
		monsters = append(monsters, &ll)
	}
	for i := 0; i < 3; i++ {
		ll := NewBulwark(ed.gamestate)
		monsters = append(monsters, &ll)
	}
	rand.Shuffle(len(monsters), func(i, j int) {
		monsters[i], monsters[j] = monsters[j], monsters[i]
	})
	ed.gamestate.AddCardToCenterDeck(DISCARD_SOURCE_NAN, true, monsters[:5]...)
}

// when slain, do this
func (ed *Firelake) OnSlain() {}

// when discarded to cooldown pile, do this
func (ed *Firelake) OnDiscarded() {}
