package cards

import "github/kharism/GuildSim_go/internal/observer"

type RemoveEventListenerAction struct {
	state    AbstractGamestate
	listener []observer.Listener
	evtName  string
}

func (r *RemoveEventListenerAction) DoAction() {
	for _, l := range r.listener {
		r.state.RemoveListener(r.evtName, l)
	}

}
func (r *RemoveEventListenerAction) AddListener(l observer.Listener) {
	r.listener = append(r.listener, l)
}
func (r *RemoveEventListenerAction) SetListener(l ...observer.Listener) {
	r.listener = l
}
func NewRemoveEventListenerAction(state AbstractGamestate, evtName string, l ...observer.Listener) AbstractActon {
	remove := &RemoveEventListenerAction{}
	remove.listener = l
	remove.state = state
	remove.evtName = evtName
	return remove
}
