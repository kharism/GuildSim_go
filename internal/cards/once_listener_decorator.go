package cards

import "github/kharism/GuildSim_go/internal/observer"

type OnceListener struct {
	listener  observer.Listener
	eventName string
	state     AbstractGamestate
}

func (o *OnceListener) DoAction(data map[string]interface{}) {
	o.listener.DoAction(data)
	o.state.RemoveListener(o.eventName, o)
}
func DecorateOnceListener(l observer.Listener, eventName string, state AbstractGamestate) observer.Listener {
	return &OnceListener{listener: l, eventName: eventName, state: state}
}
