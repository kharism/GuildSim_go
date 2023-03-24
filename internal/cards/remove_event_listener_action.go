package cards

import "github/kharism/GuildSim_go/internal/observer"

type RemoveEventListenerAction struct {
	state    AbstractGamestate
	listener observer.Listener
	evtName  string
}

func (r *RemoveEventListenerAction) DoAction() {
	r.state.RemoveListener(r.evtName, r.listener)
}
func (r *RemoveEventListenerAction) SetListener(l observer.Listener) {
	r.listener = l
}
func NewRemoveEventListenerAction(state AbstractGamestate, evtName string, l observer.Listener) AbstractActon {
	remove := &RemoveEventListenerAction{}
	remove.listener = l
	remove.state = state
	remove.evtName = evtName
	return remove
}
