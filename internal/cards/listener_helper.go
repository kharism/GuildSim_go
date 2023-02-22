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
