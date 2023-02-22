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
		t.Log("Event listener not detached")
		t.FailNow()
	}
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
