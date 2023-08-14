package cards

type GoblinSmallLairArea struct {
	BaseArea
	gamestate AbstractGamestate
}

func NewGoblinSmallLairArea(state AbstractGamestate) GoblinSmallLairArea {
	h := GoblinSmallLairArea{gamestate: state}
	return h
}

func (ed *GoblinSmallLairArea) GetName() string {
	return "GoblinSmallLairArea"
}
func (ed *GoblinSmallLairArea) GetDescription() string {
	return "Reward: 100Money and 2 Reputation also shuffle goblinwolfraider into center deck"
}
func (ed *GoblinSmallLairArea) GetCost() Cost {
	c := NewCost()
	c.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return c
}

// when played on hand, to this
func (ed *GoblinSmallLairArea) OnPlay() {}
func (ed *GoblinSmallLairArea) OnExplored() {
	// ed.gamestate.AddResource(RESOURCE_NAME_MONEY, 100)
	newPotion := ed.gamestate.GenerateRandomPotion(RARITY_COMMON)
	ed.gamestate.AddItem(newPotion)
	ed.gamestate.AddResource(RESOURCE_NAME_REPUTATION, 2)
	wolfRaiders := []Card{}
	for i := 0; i < 2; i++ {
		ll := NewGoblinWolfRaiderMonster(ed.gamestate)
		wolfRaiders = append(wolfRaiders, &ll)
	}
	ed.gamestate.AddCardToCenterDeck(DISCARD_SOURCE_NAN, true, wolfRaiders...)
}

// when slain, do this
func (ed *GoblinSmallLairArea) OnSlain() {}

// when discarded to cooldown pile, do this
func (ed *GoblinSmallLairArea) OnDiscarded() {}
