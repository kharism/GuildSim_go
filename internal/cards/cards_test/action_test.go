package cards_test

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"testing"
)

type DummyAction struct {
}

func (d *DummyAction) DoAction() {
	fmt.Println("DummyAction1")
}

type DummyAction2 struct {
}

func (d *DummyAction2) DoAction() {
	fmt.Println("DummyAction2")
}

func TestCompositeAction(t *testing.T) {
	comp := cards.NewCompositeAction(nil, &DummyAction{}, &DummyAction2{})
	comp.DoAction()
}

func TestAddResourceAction(t *testing.T) {
	state := NewDummyGamestate()
	action := cards.NewAddResourceAction(state, cards.RESOURCE_NAME_EXPLORATION, 1)
	action.DoAction()
	aa := (state.GetCurrentResource()).Detail[cards.RESOURCE_NAME_EXPLORATION]
	if aa != 1 {
		t.Log(aa)
		t.FailNow()
	}
}
