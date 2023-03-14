package cards

import "github/kharism/GuildSim_go/internal/observer"

type CardDefeatedListener struct {
	filter CardFilter
	action AbstractActon
}

func NewCardDefeatedListener(f CardFilter, action AbstractActon) observer.Listener {
	p := &CardPlayedListener{filter: f, action: action}
	return p
}
func (l *CardDefeatedListener) DoAction(data map[string]interface{}) {
	cardPlayed := data[EVENT_ATTR_CARD_PLAYED].(Card)
	if Match(cardPlayed, l.filter) {
		l.action.DoAction()
	}
}
