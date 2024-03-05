package cards

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/observer"
)

type RogueInfiltrator struct {
	BaseHero
	gamestate           AbstractGamestate
	preventTrapListener observer.Listener
}

func (r *RogueInfiltrator) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func (r *RogueInfiltrator) OnDiscarded() {
	r.gamestate.RemoveListener(EVENT_CARD_PLAYED, r.preventTrapListener)
}
func NewRogueInfiltrator(gamestate AbstractGamestate) RogueInfiltrator {
	return RogueInfiltrator{gamestate: gamestate}
}

func (r *RogueInfiltrator) GetName() string {
	return "RogueInfiltrator"
}
func (r *RogueInfiltrator) GetDescription() string {
	return "gain 3 Exploration or preemptive disarm a trap"
}
func (r *RogueInfiltrator) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return cost
}

type TrapRemovalListener struct {
	state AbstractGamestate
}

func (t *TrapRemovalListener) DoAction(data map[string]interface{}) {
	cardPlayed := data[EVENT_ATTR_CARD_DRAWN].(Card)
	fmt.Println("Detect center draw")
	if _, ok := cardPlayed.(Trapper); ok {
		j := cardPlayed.(Trapper)
		if !j.IsDisarmed() {
			j.Disarm()
			t.state.RemoveListener(EVENT_CARD_DRAWN_CENTER, t)
		}
	}
}
func (r *RogueInfiltrator) OnPlay() {
	if r.gamestate.GetBoolPicker().BoolPick("Gain 3 exploration?") {
		r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	} else {
		//r.gamestate.AddResource(RESOURCE_NAME_BLOCK, 5)
		if r.preventTrapListener == nil {
			r.preventTrapListener = &TrapRemovalListener{state: r.gamestate}
		}
		r.gamestate.AttachListener(EVENT_CARD_DRAWN_CENTER, r.preventTrapListener)
	}

}
