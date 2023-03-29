package cards_test

import (
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/cards/item"
	"github/kharism/GuildSim_go/internal/factory"
	"testing"
)

func TestPotions(t *testing.T) {
	gamestate := NewDummyGamestate()
	cardPicker := TestCardPicker{}
	cardPicker.ChooseMethod = StaticCardPicker(0)

	dumGamestate := gamestate.(*DummyGamestate)
	gamestate.SetCardPicker(&cardPicker)

	dumGamestate.CardsInDeck = cards.DeterministicDeck{}
	startHp := gamestate.GetCurrentHP()
	k := item.NewHealingPotion(gamestate)
	gamestate.AddItem(&k)
	gamestate.ConsumeItem(&k)
	nowHp := gamestate.GetCurrentHP()
	if nowHp != startHp+5 {
		t.Log("Fail", startHp, nowHp)
		t.FailNow()
	}
	if len(dumGamestate.ItemCards) > 0 {
		t.Log("Fail to remove item")
		t.FailNow()
	}
	hh := item.NewCombatPotion(gamestate)
	gamestate.AddItem(&hh)
	gamestate.ConsumeItem(&hh)
	if len(dumGamestate.ItemCards) > 0 {
		t.Log("Fail to remove item")
		t.FailNow()
	}
	if gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_COMBAT] != 3 {
		t.Log("Fail to generate resource")
		t.FailNow()
	}
	gg := item.NewExplorePotion(gamestate)
	gamestate.AddItem(&gg)
	gamestate.ConsumeItem(&gg)
	if len(dumGamestate.ItemCards) > 0 {
		t.Log("Fail to remove item")
		t.FailNow()
	}
	if gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION] != 3 {
		t.Log("Fail to generate resource")
		t.FailNow()
	}
	baseHero := cards.NewRookieAdventurer(gamestate)
	dumGamestate.CardsDiscarded.Push(&baseHero)
	aa := item.NewBanishPotion(gamestate)
	gamestate.AddItem(&aa)
	gamestate.ConsumeItem(&aa)
	if len(dumGamestate.ItemCards) > 0 {
		t.Log("Fail to remove item")
		t.FailNow()
	}
	if len(dumGamestate.ItemCards) > 0 {
		t.Log("Fail to remove item")
		t.FailNow()
	}
	if dumGamestate.CardsDiscarded.Size() > 0 {
		t.Log("Fail to banish item")
		t.FailNow()
	}
}
func TestRefreshPotion(t *testing.T) {
	gamestate := NewDummyGamestate()
	starterDeck := factory.CardFactory(factory.SET_STARTER_DECK, gamestate)
	dumGamestate := gamestate.(*DummyGamestate)
	dumGamestate.CardsInDeck.SetList(starterDeck)
	gamestate.BeginTurn() // cards in hand should be 5
	hand := gamestate.GetCardInHand()
	gamestate.PlayCard(hand[0]) // now card in hand should be 4
	rookieMage := cards.NewRookieMage(gamestate)
	dumGamestate.CardsInDeck.Stack(&rookieMage)
	handNew := gamestate.GetCardInHand()
	if len(hand) == len(handNew) {
		t.Log("Failed to remove cards")
		t.FailNow()
	}
	firstName := handNew[0].GetName()
	refre := item.NewRefreshPotion(gamestate)
	gamestate.AddItem(&refre)
	gamestate.ConsumeItem(&refre)
	handNew2 := gamestate.GetCardInHand()
	if len(handNew2) != len(handNew) {
		t.Log("Failed to draw cards")
		t.FailNow()
	}
	// t.Log(handNew2)

	if firstName == handNew2[0].GetName() {
		t.Log("fail to refresh")
		t.FailNow()
	}
}
func TestTalisman(t *testing.T) {
	gamestate := NewDummyGamestate()
	starterDeck := factory.CardFactory(factory.SET_STARTER_DECK, gamestate)
	dumGamestate := gamestate.(*DummyGamestate)
	dumGamestate.CardsInDeck.SetList(starterDeck)
	combatTalistman := item.NewCombatGauntlet(gamestate)
	explorerBoots := item.NewExplorerBoots(gamestate)
	gamestate.AddItem(&combatTalistman)
	gamestate.AddItem(&explorerBoots)
	gamestate.BeginTurn()
	if gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_COMBAT] != 2 {
		t.Log("failed to generate resource")
		t.FailNow()
	}
	if gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION] != 2 {
		t.Log("failed to generate resource")
		t.FailNow()
	}

}
