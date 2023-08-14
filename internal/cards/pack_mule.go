package cards

import "github/kharism/GuildSim_go/internal/observer"

type PackMule struct {
	RookieAdventurer
	eventListener observer.Listener
}

func NewPackMule(state AbstractGamestate) PackMule {
	base := NewRookieAdventurer(state)
	this := PackMule{base, nil}
	return this
}

func (r *PackMule) GetName() string {
	return "Pack mule"
}
func (r *PackMule) GetDescription() string {
	return "Add 1 exlporation point, gains additional 1 Exploration point if you play or have played advanced adventurer this turn"
}
func (r *PackMule) GetCost() Cost {
	cost := NewCost()
	cost.Resource.Detail[RESOURCE_NAME_MONEY] = 50
	return cost
}

type AddResourceOnce struct {
	state        AbstractGamestate
	resourceName string
	cardEvent    observer.Listener
	amount       int
}

func (a *AddResourceOnce) DoAction() {
	a.state.AddResource(a.resourceName, a.amount)
}

func (r *PackMule) OnDiscarded() {
	r.gamestate.RemoveListener(EVENT_CARD_PLAYED, r.eventListener)
}
func (r *PackMule) Dispose(source string) {
	r.gamestate.DiscardCard(r, source)
}
func (r *PackMule) OnPlay() {
	r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 1)
	aa := AdvancedAdventurer{}
	filter := &CardFilter{Key: FILTER_NAME, Op: Eq, Value: aa.GetName()}
	if Contains(r.gamestate.GetPlayedCards(), filter) {
		r.gamestate.AddResource(RESOURCE_NAME_EXPLORATION, 1)
	} else {
		// add one time event listener
		if r.eventListener == nil {
			addResourceAction := NewAddResourceAction(r.gamestate, RESOURCE_NAME_EXPLORATION, 1)
			removeEventListenerAction := NewRemoveEventListenerAction(r.gamestate, EVENT_CARD_PLAYED, nil)
			compositeAction := NewCompositeAction(r.gamestate, addResourceAction, removeEventListenerAction)
			cardPlayedListener := NewCardPlayedListener(filter, compositeAction)
			removeEventListenerAction.(*RemoveEventListenerAction).listener = append(removeEventListenerAction.(*RemoveEventListenerAction).listener, cardPlayedListener)
			r.eventListener = cardPlayedListener
		}

		// fmt.Println(removeEventListenerAction)
		// packMuleListener := DecorateOnceListener(CardPlayedListener, EVENT_ATTR_CARD_PLAYED, r.gamestate)
		r.gamestate.AttachListener(EVENT_CARD_PLAYED, r.eventListener)
	}
}
