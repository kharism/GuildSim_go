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

func TestShockCurse(t *testing.T) {
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
	shockCurse := cards.NewShockCurse(gamestate)

	dumGamestate.StackCards("NA", &shockCurse)
	t.Log(dumGamestate.RuleEnforcer[cards.ACTION_DRAW])
	dumGamestate.BeginTurn()
	if dumGamestate.RuleEnforcer[cards.ACTION_DRAW].Len() != 1 {
		t.Log("Failed to attach enforcer")
		t.Fail()
	}
	dumGamestate.StackCards(cards.DISCARD_SOURCE_HAND, &shockCurse)
	if dumGamestate.RuleEnforcer[cards.ACTION_DRAW].Len() != 0 {
		t.Log("Failed to detach enforcer", dumGamestate.RuleEnforcer[cards.ACTION_DRAW].Len())
		t.Fail()
	}
}
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
func TestWolfShaman(t *testing.T) {
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

	wolfShaman := cards.NewWolfShaman(gamestate)
	gamestate.PlayCard(&wolfShaman)
	iceWyvern := cards.NewIceWyvern(gamestate)
	freezeCurse := cards.NewFreezeCurse(gamestate)
	dumGamestate.CenterCards = append(dumGamestate.CenterCards, &iceWyvern)
	dumGamestate.EndTurn()
	// t.Log(dumGamestate.CardsInDeck.Deck.List()[0].GetName())
	if dumGamestate.CardsInDeck.Deck.List()[0].GetName() == freezeCurse.GetName() {
		t.FailNow()
	}
	dumGamestate.EndTurn()
	// t.Log(dumGamestate.CardsInDeck.Deck.List()[0].GetName())
	if dumGamestate.CardsInDeck.Deck.List()[0].GetName() != freezeCurse.GetName() {
		t.FailNow()
	}
	dumGamestate.Draw()
	dumGamestate.CenterCards = append(dumGamestate.CenterCards, &iceWyvern)
	gamestate.PlayCard(&wolfShaman)
	dumGamestate.EndTurn()
	if dumGamestate.CardsInDeck.Deck.List()[0].GetName() != freezeCurse.GetName() {
		t.FailNow()
	}
}
