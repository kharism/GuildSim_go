package cards_test

import (
	"github/kharism/GuildSim_go/internal/cards"
	"testing"
)

func TestHealingPotion(t *testing.T) {
	gamestate := NewDummyGamestate()
	cardPicker := TestCardPicker{}

	cardPicker.ChooseMethod = StaticCardPicker(1)

	dumGamestate := gamestate.(*DummyGamestate)

	dumGamestate.CardsInDeck = cards.DeterministicDeck{}
	dumGamestate.cardPiker = &cardPicker

	k := cards.NewRookieMage(gamestate)
	l := cards.NewStagShaman(gamestate)
	dumGamestate.CardsInDeck.Push(&k)
	dumGamestate.CardsInDeck.Push(&l)

	for i := 0; i < 10; i++ {
		j := cards.NewRookieAdventurer(gamestate)
		dumGamestate.CardsInDeck.Push(&j)
	}

	for i := 0; i < 5; i++ {
		gamestate.Draw()
	}
	cardsInHand := gamestate.GetCardInHand()
	if len(cardsInHand) != 5 {
		t.Log("Failed to draw")
		t.FailNow()
	}
	gamestate.PlayCard(cardsInHand[0])
	if len(cardsInHand) != 5 {
		t.Log("Failed to draw")
		t.FailNow()
	}

	if dumGamestate.CardsDiscarded.Size() != 1 {
		t.Log("Failed to discard")
		t.FailNow()
	}
	cardPicker.ChooseMethod = StaticCardPicker(0)
	gamestate.PlayCard(&l)

	if dumGamestate.CardsDiscarded.Size() != 0 {
		t.Log("Failed to banish from cooldown pile")
		t.FailNow()
	}
	if len(cardsInHand) != 5 {
		t.Log("Failed to draw")
		t.Fail()
	}

}
