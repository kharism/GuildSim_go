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
	d.state.AddCardToCenterDeck(&newHero)
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
