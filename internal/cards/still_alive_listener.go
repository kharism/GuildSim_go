package cards

import "github/kharism/GuildSim_go/internal/observer"

type StillAliveListener struct {
	state  AbstractGamestate
	action AbstractActon
}

// this listener must be added to gamestate to ensure it knows when the player die
// will execute action when HP reaches 0 or lower
func NewStillAliveListener(state AbstractGamestate, action AbstractActon) observer.Listener {
	return &StillAliveListener{state: state, action: action}
}

func (s *StillAliveListener) DoAction(data map[string]interface{}) {
	if s.state.GetCurrentHP() <= 0 {
		s.action.DoAction()
	}
}
