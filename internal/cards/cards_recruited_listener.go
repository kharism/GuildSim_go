package cards

import "github/kharism/GuildSim_go/internal/observer"

type CardRecruitedListener struct {
	filter *CardFilter
	action AbstractActon
}

func NewCardRecruitedListener(f *CardFilter, action AbstractActon) observer.Listener {
	p := &CardRecruitedListener{filter: f, action: action}
	return p
}
func (l *CardRecruitedListener) DoAction(data map[string]interface{}) {
	cardRecruited := data[EVENT_ATTR_CARD_RECRUITED].(Card)
	if l.filter != nil && Match(cardRecruited, l.filter) {
		l.action.DoAction()
	} else {
		l.action.DoAction()
	}
}
