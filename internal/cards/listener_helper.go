package cards

import (
	"github/kharism/GuildSim_go/internal/observer"
)

type OneTimeListener struct {
	l  observer.Listener
	ob observer.Observer
}

func (o *OneTimeListener) DoAction(data map[string]interface{}) {
	o.l.DoAction(data)
	o.ob.Detach(o)
}
func MakeOneTimeListener(l observer.Listener, ob observer.Observer) observer.Listener {
	o := &OneTimeListener{l: l, ob: ob}
	return o
}

// a simple listener that execute action when triggered
// it does not check the data sent by sender
type BasicListener struct {
	action AbstractActon
}

func NewBasicAction(action AbstractActon) *BasicListener {
	return &BasicListener{action: action}
}
func (b *BasicListener) DoAction(data map[string]interface{}) {
	b.action.DoAction()
}
