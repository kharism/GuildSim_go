package cards

import "github/kharism/GuildSim_go/internal/observer"

type CardPlayedListener struct {
	filter CardFilter
	action AbstractActon
}

func NewCardPlayedListener(f CardFilter, action AbstractActon) observer.Listener {
	p := &CardPlayedListener{filter: f, action: action}
	return p
}
func (l *CardPlayedListener) DoAction(data map[string]interface{}) {
	cardPlayed := data[EVENT_ATTR_CARD_PLAYED].(Card)
	if Match(cardPlayed, l.filter) {
		l.action.DoAction()
	}
}
