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

	dgs := gamestate.(*DummyGamestate)
	rokieCombatant := cards.NewRookieCombatant(gamestate)
	dgs.CardsDiscarded.Push(&rokieCombatant)
	dgs.PlayCard(&rookieNurse)
	if len(dgs.CardsInHand) != 1 {
		t.Log("failed to draw")
		t.FailNow()
	}
	if dgs.CardsInDeck.Size() != 0 {
		t.Log("failed to shuffle back 1")
		t.FailNow()
	}
	if dgs.CardsDiscarded.Size() != 0 {
		t.Log("failed to shuffle back 2")
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
	centerDeck := cards.Deck{}
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
