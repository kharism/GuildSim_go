package cards_test

import (
	"github/kharism/GuildSim_go/internal/cards"
	"math/rand"
	"testing"
)

func PrintCardList(list []cards.Card, t *testing.T) {
	for idx, i := range list {
		t.Log(idx, i.GetName(), i.GetCost().Detail, i.GetDescription())
	}
}
func TestTombEntrance(t *testing.T) {
	rand.Seed(12)
	gamestate := NewDummyGamestate()

	tombEntrance := cards.NewTombMonarchEntrance(gamestate)

	dumGamestate := gamestate.(*DummyGamestate)

	dumGamestate.CenterCards = append(dumGamestate.CenterCards, &tombEntrance)
	starterDeck := []cards.Card{}
	for i := 0; i < 5; i++ {
		newAdventurer := cards.NewRookieAdventurer(gamestate)
		newCombatant := cards.NewRookieCombatant(gamestate)
		starterDeck = append(starterDeck, &newAdventurer, &newCombatant)
	}
	dumGamestate.CardsInDeck.SetList(starterDeck)
	for i := 0; i < 5; i++ {
		gamestate.Draw()
	}
	cardInHand := gamestate.GetCardInHand()
	if len(cardInHand) != 5 {
		t.Log("failed to draw")
		t.FailNow()
	}

	// only play adventurer card
	gamestate.PlayCard(cardInHand[0])
	gamestate.PlayCard(cardInHand[1])
	gamestate.PlayCard(cardInHand[2])
	curRes := gamestate.GetCurrentResource()
	if curRes.Detail[cards.RESOURCE_NAME_EXPLORATION] != 3 {
		t.Log("Exploration points is wrong")
		t.FailNow()
	}
	gamestate.Explore(&tombEntrance)
	if len(dumGamestate.TopicsListeners[cards.EVENT_CARD_DEFEATED].Listeners) != 1 {
		t.Log("failed to attach listener")
		t.FailNow()
	}
	t.Log(dumGamestate.CardsInCenterDeck.Size())
	PrintCardList(dumGamestate.CenterCards, t)
	t.Log(gamestate.GetCurrentHP())
	gamestate.PlayCard(cardInHand[0])
	gamestate.PlayCard(cardInHand[0])
	gamestate.EndTurn()
	t.Log("===END===")
	if gamestate.GetCurrentHP() != 57 {
		t.Log("failed to take damage")
		t.FailNow()
	}
	for i := 0; i < 5; i++ {
		gamestate.Draw()
	}
	PrintCardList(dumGamestate.CardsInHand, t)
	cardInHand = dumGamestate.CardsInHand
	gamestate.PlayCard(cardInHand[4])
	gamestate.PlayCard(cardInHand[2])
	gamestate.PlayCard(cardInHand[0])
	// t.Log(dumGamestate.CenterCards[0].GetName(), dumGamestate.CardsInCenterDeck.Size(), gamestate.GetCurrentResource().Detail)
	// PrintCardList(dumGamestate.CenterCards, t)
	t.Log(len(dumGamestate.CenterCards), dumGamestate.CenterCards[0].GetName())
	gamestate.DefeatCard(dumGamestate.CenterCards[0])
	t.Log(dumGamestate.CardsInCenterDeck.Size())
	if dumGamestate.CardsInCenterDeck.Size() != 1 {
		t.Log("failed to defeat and replace")
		t.FailNow()
	}
	gamestate.EndTurn()
	t.Log("===END===")
	for i := 0; i < 8; i++ {
		gamestate.Draw()
	}
	PrintCardList(dumGamestate.CenterCards, t)
	PrintCardList(dumGamestate.CardsInHand, t)
	cardInHand = dumGamestate.CardsInHand
	gamestate.PlayCard(cardInHand[6])
	gamestate.PlayCard(cardInHand[4])
	gamestate.PlayCard(cardInHand[3])
	gamestate.PlayCard(cardInHand[2])
	gamestate.PlayCard(cardInHand[1])
	cardPicker := TestCardPicker{ChooseMethod: StaticCardPicker(0)}
	gamestate.SetCardPicker(&cardPicker)
	gamestate.DefeatCard(dumGamestate.CenterCards[0])
	t.Log(dumGamestate.CardsInCenterDeck.Size())
	if dumGamestate.CardsInCenterDeck.Size() != 1 {
		PrintCardList(dumGamestate.CardsInCenterDeck.List(), t)
		t.Log("failed to defeat and replace")
		t.FailNow()
	}

}
