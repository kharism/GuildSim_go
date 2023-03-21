package cards_test

import (
	"github/kharism/GuildSim_go/internal/cards"
	"testing"
)

func TestRookieAdventurer(t *testing.T) {
	gamestate := NewDummyGamestate()
	rookieAdventurer := cards.NewRookieAdventurer(gamestate)
	//rookieAdventurer.OnPlay()
	advancedAdventurer := cards.NewAdvancedAdventurer(gamestate)
	packMule := cards.NewPackMule(gamestate)
	gamestate.PlayCard(&rookieAdventurer)
	if _, ok := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]; ok {
		hh := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]
		if hh != 1 {
			t.FailNow()
		}
	}
	gamestate.PlayCard(&advancedAdventurer)
	if _, ok := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]; ok {
		hh := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]
		if hh != 3 {
			t.FailNow()
		}
	}
	gamestate.PlayCard(&packMule)
	if _, ok := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]; ok {
		hh := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]
		if hh != 5 {
			t.FailNow()
		}
	}
	if len(gamestate.GetPlayedCards()) != 3 {
		t.FailNow()
	}
}
func TestNurse(t *testing.T) {
	gamestate := NewDummyGamestate()

	rookieNurse := cards.NewRookieNurse(gamestate)
	hpStart := gamestate.GetCurrentHP()
	dgs := gamestate.(*DummyGamestate)
	rokieCombatant := cards.NewRookieCombatant(gamestate)
	dgs.CardsDiscarded.Push(&rokieCombatant)
	dgs.PlayCard(&rookieNurse)
	hpCurrent := gamestate.GetCurrentHP()
	if hpCurrent-hpStart != 1 {
		t.Log("failed to recover")
		t.FailNow()
	}
}
func TestPackMuleEventListener(t *testing.T) {
	gamestate := NewDummyGamestate()
	advancedAdventurer := cards.NewAdvancedAdventurer(gamestate)
	packMule := cards.NewPackMule(gamestate)
	gamestate.PlayCard(&packMule)
	if _, ok := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]; ok {
		hh := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]
		if hh != 1 {
			t.FailNow()
		}
	}
	evtListerner := gamestate.(*DummyGamestate).TopicsListeners[cards.EVENT_CARD_PLAYED]
	if len(evtListerner.Listeners) != 1 {
		t.Log("Event listener not attached")
		t.FailNow()
	}
	gamestate.EndTurn()
	evtListerner = gamestate.(*DummyGamestate).TopicsListeners[cards.EVENT_CARD_PLAYED]
	if len(evtListerner.Listeners) != 0 {
		t.Log("Event listener not detached")
		t.FailNow()
	}
	gamestate.PlayCard(&packMule)
	gamestate.PlayCard(&advancedAdventurer)
	if _, ok := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]; ok {
		hh := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]
		if hh != 4 {
			t.Log(hh)
			t.FailNow()
		}
	}
	evtListerner = gamestate.(*DummyGamestate).TopicsListeners[cards.EVENT_CARD_PLAYED]
	if len(evtListerner.Listeners) != 0 {
		t.FailNow()
	}
	gamestate.PlayCard(&advancedAdventurer)
	if _, ok := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]; ok {
		hh := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION]
		if hh != 6 {
			t.Log(hh)
			t.FailNow()
		}
	}
}

func TestRookieCombatant(t *testing.T) {
	gamestate := NewDummyGamestate()
	rookieCombatant := cards.NewRookieCombatant(gamestate)
	gamestate.PlayCard(&rookieCombatant)
	if _, ok := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_COMBAT]; ok {
		hh := gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_COMBAT]
		if hh != 1 {
			t.FailNow()
		}
	}
}

func TestExpore(t *testing.T) {
	gamestate := NewDummyGamestate()
	dungeon1 := cards.NewEasyDungeonArea(gamestate)
	baseHero1 := cards.BaseHero{}
	centerDeck := cards.DeterministicDeck{}
	centerDeck.Push(&dungeon1)
	centerDeck.Push(&baseHero1)
	gamestate.(*DummyGamestate).CardsInCenterDeck = centerDeck
	gamestate.CenterRowInit()
	if len(gamestate.GetCenterCard()) != 1 {
		t.Error("DDDDD")
		t.FailNow()
	}

	rookieAdv := cards.NewRookieAdventurer(gamestate)
	advAdv := cards.NewAdvancedAdventurer(gamestate)

	gamestate.PlayCard(&rookieAdv)
	gamestate.PlayCard(&advAdv)

	gamestate.Explore(&dungeon1)

	// check rewards
	if gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_MONEY] != 100 {
		t.Error("Failed to gain resource")
		t.FailNow()
	}

	// check replacement
	centerCard := gamestate.GetCenterCard()
	if len(centerCard) != 1 {
		t.Error("Failed to Replace card in center row")
		t.FailNow()
	}
	if centerCard[0].GetName() == "EasyDungeonArea" {
		t.Error("Failed to Replace card in center row")
		t.FailNow()
	}

	// check current resource
	if gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_EXPLORATION] != 0 {
		t.Error("Resource is not reduced")
		t.FailNow()
	}
}

func TestDefeat(t *testing.T) {
	gamestate := NewDummyGamestate()
	monster := cards.BaseMonster{}
	goblin := cards.NewGoblinMonster(gamestate)

	combatant := cards.NewRookieCombatant(gamestate)
	gamestate.PlayCard(&combatant)
	gamestate.(*DummyGamestate).CardsInCenterDeck.Push(&monster)
	gamestate.(*DummyGamestate).CenterCards = append(gamestate.(*DummyGamestate).CenterCards, &goblin)
	gamestate.DefeatCard(&goblin)

	if gamestate.GetCurrentResource().Detail[cards.RESOURCE_NAME_COMBAT] != 0 {
		t.Error("Resource is not reduced")
		t.FailNow()
	}

	// check replacement
	centerCard := gamestate.GetCenterCard()
	if len(centerCard) != 1 {
		t.Error("Failed to Replace card in center row")
		t.FailNow()
	}
	if centerCard[0].GetName() == "GoblinMonster" {
		t.Error("Failed to Replace card in center row")
		t.FailNow()
	}

}

func TestPunish(t *testing.T) {
	gamestate := NewDummyGamestate()
	monster := cards.BaseMonster{}
	goblinWolfRaider := cards.NewGoblinWolfRaiderMonster(gamestate)
	hpStart := gamestate.GetCurrentHP()
	gamestate.(*DummyGamestate).CardsInCenterDeck.Push(&monster)
	gamestate.(*DummyGamestate).CenterCards = append(gamestate.(*DummyGamestate).CenterCards, &goblinWolfRaider)

	gamestate.EndTurn()
	hpAfterPunish := gamestate.GetCurrentHP()
	t.Log(hpAfterPunish)
	if hpAfterPunish != hpStart-2 {
		t.Log("Failed to inflict punish")
		t.FailNow()
	}
}

func TestEndTurn(t *testing.T) {
	gamestate := NewDummyGamestate()
	monster := cards.BaseMonster{}
	goblinWolfRaider := cards.NewGoblinWolfRaiderMonster(gamestate)
	hpStart := gamestate.GetCurrentHP()
	dummyGamestate := gamestate.(*DummyGamestate)
	dummyGamestate.CardsInCenterDeck.Push(&monster)
	dummyGamestate.CenterCards = append(gamestate.(*DummyGamestate).CenterCards, &goblinWolfRaider)

	adventurer1 := cards.NewRookieAdventurer(gamestate)
	adventurer2 := cards.NewRookieAdventurer(gamestate)

	baseHero1 := cards.BaseHero{}
	onDrawCurse := cards.NewDamageDrawCurse(gamestate)
	onEndTurnCurse := cards.NewDamageEndturnCurse(gamestate)

	dummyGamestate.CardsInHand = []cards.Card{&adventurer1, &adventurer2}
	dummyGamestate.CardsInDeck = cards.DeterministicDeck{}

	dummyGamestate.CardsInDeck.Push(&onDrawCurse)
	dummyGamestate.CardsInDeck.Push(&baseHero1)
	dummyGamestate.CardsInDeck.Push(&onEndTurnCurse)
	// wolf raider inflct damage here
	gamestate.EndTurn()
	if dummyGamestate.CardsDiscarded.Size() != 2 {
		t.Log("It should be 2")
		t.FailNow()
	}
	// draw 3 cards
	for i := 0; i < 3; i++ {
		gamestate.Draw()
	}

	if len(dummyGamestate.CardsInHand) != 2 {
		t.Log("cards in hand should be 2")
		t.FailNow()
	}
	if len(dummyGamestate.CardsBanished) != 1 {
		t.Log("It should be 1 , but it is", len(dummyGamestate.CardsBanished), " instead")
		t.FailNow()
	}
	hpNow := gamestate.GetCurrentHP()
	// 4 dmg from wolfraider+damage draw curse
	if hpStart-hpNow != 4 {
		t.Log("Failed to inflict 4 dmg ", hpStart-hpNow)
		t.FailNow()
	}
	// inflict damage on end turn by wolf raider+end turn curse
	gamestate.EndTurn()
	hpNow = gamestate.GetCurrentHP()
	if hpStart-hpNow != 8 {
		t.Log("Failed to inflict 8 dmg ", hpStart-hpNow)
		t.FailNow()
	}
}

func TestRookieMage(t *testing.T) {
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

func TestWingedLion(t *testing.T) {
	gamestate := NewDummyGamestate()
	cardPicker := TestCardPicker{}

	cardPicker.ChooseMethod = StaticCardPicker(1)
	dumGamestate := gamestate.(*DummyGamestate)
	dumGamestate.cardPiker = &cardPicker

	area1 := cards.NewGoblinSmallLairArea(gamestate)
	area2 := cards.NewEasyDungeonArea(gamestate)
	monster1 := cards.NewGoblinMonster(gamestate)
	monster2 := cards.NewGoblinWolfRaiderMonster(gamestate)
	monster3 := cards.NewGoblinWolfRaiderMonster(gamestate)
	dumGamestate.CenterCards = append(dumGamestate.CenterCards, &area1, &area2, &monster1, &monster2)

	dumGamestate.CardsInCenterDeck.Push(&monster3)
	for i := 0; i < 10; i++ {
		j := cards.NewRookieMage(gamestate)
		dumGamestate.CardsInCenterDeck.Push(&j)
		k := cards.NewRookieAdventurer(gamestate)
		dumGamestate.CardsInDeck.Push(&k)
	}
	wingedLion := cards.NewWingedLion(gamestate)
	dumGamestate.CardsInHand = append(dumGamestate.CardsInHand, &wingedLion)
	gamestate.PlayCard(&wingedLion)
	newCenterCard := dumGamestate.CenterCards
	// t.Log(newCenterCard[1].GetName())
	if newCenterCard[1].GetName() == area2.GetName() {
		t.Log("Failed to replace")
		t.FailNow()
	}
	cardInHand := gamestate.GetCardInHand()
	// t.Log(cardInHand[0].GetName())
	if len(cardInHand) != 1 {
		t.Log("Failed to Draw")
		t.FailNow()
	}

}
