package cards_test

import (
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/observer"
	"math/rand"
	"testing"
)

type DummyEventListener struct {
	Listeners []observer.Listener
}

func NewDummyEventListener() DummyEventListener {
	d := DummyEventListener{}
	d.Listeners = []observer.Listener{}
	return d
}
func (d *DummyEventListener) Attach(l observer.Listener) {

	d.Listeners = append(d.Listeners, l)
	// fmt.Println("attach from del", len(d.Listeners))
}

func (d *DummyEventListener) Detach(l observer.Listener) {
	newList := []observer.Listener{}
	for _, i := range d.Listeners {
		if i == l {
			continue
		} else {
			newList = append(newList, i)
		}
	}
	d.Listeners = newList
}

func (d *DummyEventListener) Notify(data map[string]interface{}) {
	for _, i := range d.Listeners {
		i.DoAction(data)
	}
}

// Dummygamestate implements abstractgamestate and publisher
type DummyGamestate struct {
	currentResource   cards.Resource
	CardsInDeck       cards.Deck
	CardsInCenterDeck cards.Deck
	TopicsListeners   map[string]*DummyEventListener
	CardsInHand       []cards.Card
	CardsPlayed       []cards.Card
	CenterCards       []cards.Card
}

func (d *DummyGamestate) AttachListener(eventName string, l observer.Listener) {
	if _, ok := d.TopicsListeners[eventName]; !ok {
		d.TopicsListeners[eventName] = &DummyEventListener{}
	}
	k := (d.TopicsListeners[eventName])
	k.Attach(l)
	// fmt.Println("Attach Listener", len(d.TopicsListeners[eventName].Listeners))
}
func (d *DummyGamestate) RemoveListener(eventName string, l observer.Listener) {
	if _, ok := d.TopicsListeners[eventName]; !ok {
		return
	}
	// fmt.Println("Remove Listener")
	k := (d.TopicsListeners[eventName])
	k.Detach(l)
}

func NewDummyGamestate() cards.AbstractGamestate {
	d := DummyGamestate{}
	d.currentResource = cards.NewResource()
	d.CardsPlayed = []cards.Card{}
	d.TopicsListeners = map[string]*DummyEventListener{}
	d.CenterCards = []cards.Card{}
	return &d
}
func (d *DummyGamestate) PlayCard(c cards.Card) {
	c.OnPlay()
	// fmt.Println("Card played", c.GetName())
	cardPlayedEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_PLAYED: c}

	l, ok := d.TopicsListeners[cards.EVENT_CARD_PLAYED]
	if ok {
		l.Notify(cardPlayedEvent)
	}

	d.CardsPlayed = append(d.CardsPlayed, c)
}
func (d *DummyGamestate) GetPlayedCards() []cards.Card {
	return d.CardsPlayed
}
func (d *DummyGamestate) GetCardInHand() []cards.Card {
	return d.CardsInHand
}
func (d *DummyGamestate) GetCenterCard() []cards.Card {
	return d.CenterCards
}
func (d *DummyGamestate) RecruitCard(c cards.Card) {
	return
}
func (d *DummyGamestate) DiscardCard(c cards.Card) {
	return
}
func (d *DummyGamestate) CenterRowInit() {
	f := d.ReplaceCenterCard()
	d.CenterCards = append(d.CenterCards, f)
}
func (d *DummyGamestate) Explore(c cards.Card) {
	// check cost and resource
	f := c.GetCost()
	res := d.currentResource
	if (&f).IsEnough(res) {
		c.OnExplored()
		cardExploredEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_EXPLORED: c}

		l, ok := d.TopicsListeners[cards.EVENT_CARD_EXPLORED]
		if ok {
			l.Notify(cardExploredEvent)
		}
		// remove c from center cards
		replacementCard := d.ReplaceCenterCard()
		newCenterCards := []cards.Card{}
		for _, v := range d.CenterCards {
			if v == c {
				newCenterCards = append(newCenterCards, replacementCard)
			} else {
				newCenterCards = append(newCenterCards, v)
			}
		}
		d.CenterCards = newCenterCards
	}
}
func (d *DummyGamestate) ReplaceCenterCard() cards.Card {
	return d.CardsInCenterDeck.Draw()
}
func (d *DummyGamestate) Draw() {
	newCard := d.CardsInDeck.Draw()
	d.CardsInHand = append(d.CardsInHand, newCard)
	return
}
func (d *DummyGamestate) BanishCard(c cards.Card) {
	return
}
func (d *DummyGamestate) DefeatCard(c cards.Card) {
	return
}
func (d *DummyGamestate) GetCurrentResource() cards.Resource {
	return d.currentResource
}
func (d *DummyGamestate) AddResource(name string, amount int) {
	d.currentResource.AddResource(name, amount)
}

func TestDeck(t *testing.T) {
	deck := cards.Deck{}
	gamestate := DummyGamestate{}
	rookieCombatant := cards.NewRookieCombatant(&gamestate)
	rookieAdventurer := cards.NewRookieCombatant(&gamestate)
	// deck = append(deck, &rookieAdventurer, &rookieCombatant, &rookieCombatant, &rookieCombatant)
	deck.Push(&rookieAdventurer)
	deck.Push(&rookieCombatant)
	deck.Push(&rookieCombatant)
	deck.Push(&rookieCombatant)
	if deck.Size() != 4 {
		t.Error("Seharusnya 4")
	}
	rand.Seed(10)

	top := deck.Draw()
	if deck.Size() != 3 {
		t.Error("Seharusnya 3")
		t.FailNow()
	}
	t.Log(top.GetName())
}
