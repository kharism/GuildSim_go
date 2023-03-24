package cards_test

import (
	"github/kharism/GuildSim_go/internal/cards"
	"testing"
)

type dummyAction struct {
	state cards.AbstractGamestate
}

func (d *dummyAction) DoAction() {
	newHero := cards.NewRookieNurse(d.state)
	d.state.AddCardToCenterDeck(cards.DISCARD_SOURCE_NAN, true, &newHero)
}
func TestCardsRecruitedListener(t *testing.T) {
	gamestate := NewDummyGamestate()
	dumbAction := &dummyAction{state: gamestate}

	dummyHero := cards.NewRookieAdventurer(gamestate)
	recruitListener := cards.NewCardRecruitedListener(nil, dumbAction)
	gamestate.AttachListener(cards.EVENT_CARD_RECRUITED, recruitListener)
	dumGs := gamestate.(*DummyGamestate)
	dumGs.CenterCards = append(dumGs.CenterCards, &dummyHero)
	gamestate.RecruitCard(&dummyHero)
	t.Log(gamestate.GetCenterCard())
	if len(gamestate.GetCenterCard()) != 1 {
		t.Log("Failed to trigger")
		t.FailNow()
	}
}

var TEST_GAMEOVER = false

type dummyAction2 struct {
}

func (d *dummyAction2) DoAction() {
	TEST_GAMEOVER = true
}

func TestStillAlive(t *testing.T) {
	gamestate := NewDummyGamestate()
	dumbAction := &dummyAction2{}
	stillAliveListener := cards.NewStillAliveListener(gamestate, dumbAction)
	gamestate.AttachListener(cards.EVENT_TAKE_DAMAGE, stillAliveListener)
	gamestate.TakeDamage(100)
	if !TEST_GAMEOVER {
		t.Log("Gagal trigger game over")
		t.FailNow()
	}
}
