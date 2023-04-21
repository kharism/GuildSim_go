package cards

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/observer"
)

type LimitDraw struct {
	state           AbstractGamestate
	drawCount       int
	limitNumber     int
	beginTurnAction observer.Listener
	endTurnAction   observer.Listener
}

func NewLimitDraw(state AbstractGamestate, limitNumber int) *LimitDraw {
	ld := &LimitDraw{state: state, limitNumber: limitNumber}
	ld.beginTurnAction = &LimitAction{parent: ld}
	ld.endTurnAction = &ReleaseLimitAction{parent: ld}
	return ld
}
func (i *LimitDraw) AttachLimitDraw(state AbstractGamestate) AbstractGamestate {
	state.AttachListener(EVENT_START_OF_TURN, i.beginTurnAction)
	state.AttachListener(EVENT_END_OF_TURN, i.endTurnAction)
	return state
}
func (i *LimitDraw) DetachLimitDraw(state AbstractGamestate) AbstractGamestate {
	fmt.Println("Reset LimitDraw")
	state.RemoveListener(EVENT_START_OF_TURN, i.beginTurnAction)
	state.RemoveListener(EVENT_END_OF_TURN, i.endTurnAction)
	state.DetachLegalCheck(ACTION_DRAW, i)
	return state
}

// func (i *LimitDraw) String() string {
// 	return "LimitDraw"
// }

// set limit after we draw for turn
type LimitAction struct {
	parent *LimitDraw
}

func (i *LimitAction) DoAction(data map[string]interface{}) {
	// i.parent.drawCount = 0
	fmt.Println("Begin turn attach limiter")
	i.parent.state.AttachLegalCheck(ACTION_DRAW, i.parent)
}

// release limit on end of turn
type ReleaseLimitAction struct {
	parent *LimitDraw
}

func (i *ReleaseLimitAction) DoAction(data map[string]interface{}) {
	i.parent.drawCount = 0
	fmt.Println("End turn reset limiter")
	i.parent.state.DetachLegalCheck(ACTION_DRAW, i.parent)
}
func (i *LimitDraw) DoAction(data map[string]interface{}) {
	i.drawCount = 0

}
func (i *LimitDraw) Check(data interface{}) bool {
	// fmt.Println("check limit", i.drawCount)
	if i.drawCount < i.limitNumber {
		i.drawCount++
		return true
	}
	return false
}
