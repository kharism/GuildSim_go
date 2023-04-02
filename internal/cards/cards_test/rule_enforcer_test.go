package cards_test

import (
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/factory"
	"testing"
)

type IllegalAction struct{}

func (i *IllegalAction) Check(data interface{}) bool {
	return false
}

var (
	NO = IllegalAction{}
)

func TestRuleEnforcer(t *testing.T) {
	gamestate := NewDummyGamestate()
	starterDeck := factory.CardFactory(factory.SET_STARTER_DECK, gamestate)
	centerCards := factory.CardFactory(factory.SET_CENTER_DECK_1, gamestate)
	dumGamestate := gamestate.(*DummyGamestate)
	dumGamestate.CardsInDeck.SetList(starterDeck)
	dumGamestate.CardsInCenterDeck.SetList(centerCards)
	cardPicker := TestCardPicker{}
	cardPicker.ChooseMethod = StaticCardPicker(0)
	cardPicker.ChooseMethodBool = func() bool { return false }
	gamestate.SetBoolPicker(&cardPicker)
	gamestate.SetCardPicker(&cardPicker)
	gamestate.SetDetailViewer(&cardPicker)

	gamestate.AttachLegalCheck(cards.ACTION_DRAW, &NO)
	gamestate.Draw()
	if len(gamestate.GetCardInHand()) != 0 {
		t.Log("fail to prevent draw")
		t.FailNow()
	}
	gamestate.DetachLegalCheck(cards.ACTION_DRAW, &NO)
	limitDraw := cards.NewLimitDraw(gamestate, 3)
	limitDraw.AttachLimitDraw(gamestate)
	gamestate.BeginTurn()
	gamestate.Draw()
	gamestate.Draw()
	gamestate.Draw()
	if len(gamestate.GetCardInHand()) != 8 {
		t.Log("fail draw", len(gamestate.GetCardInHand()))
		t.FailNow()
	}
	gamestate.Draw()
	if len(gamestate.GetCardInHand()) != 8 {
		t.Log("fail preventing draw", len(gamestate.GetCardInHand()))
		t.FailNow()
	}
	gamestate.EndTurn()
	gamestate.BeginTurn()
	// t.Log(len(gamestate.GetCardInHand()))
	if len(gamestate.GetCardInHand()) != 5 {
		t.Log("fail release limiter", len(gamestate.GetCardInHand()))
		t.FailNow()
	}
	gamestate.Draw()
	gamestate.Draw()
	gamestate.Draw()
	if len(gamestate.GetCardInHand()) != 8 {
		t.Log("fail draw", len(gamestate.GetCardInHand()))
		t.FailNow()
	}
}
