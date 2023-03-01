package cards_test

import (
	"bufio"
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/observer"
	"math/rand"
	"os"
	"strconv"
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

// show a basic card picker using stdin/stdout
type TextCardPicker struct{}

func (t *TextCardPicker) PickCard(list []cards.Card, message string) int {
	fmt.Println(message)
	for i, card := range list {
		fmt.Printf("[%d] %s [%s]\n", i, card.GetName(), card.GetCost())
	}
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	for {
		picks, err := strconv.Atoi(text)
		if err != nil {
			continue
		}
		return picks
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
	HitPoint          int
	//ui stuff
	cardPiker cards.AbstractCardPicker
}

func (d *DummyGamestate) PayResource(cost cards.Cost) {
	for key, val := range cost.Detail {
		d.currentResource.Detail[key] -= val
	}
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
	d.CardsInHand = []cards.Card{}
	d.cardPiker = &TextCardPicker{}
	d.HitPoint = 60
	return &d
}
func (d *DummyGamestate) GetCurrentHP() int {
	return d.HitPoint
}
func (d *DummyGamestate) TakeDamage(dmg int) {
	d.HitPoint -= dmg
	l, ok := d.TopicsListeners[cards.EVENT_TAKE_DAMAGE]
	takeDamageEvent := map[string]interface{}{cards.EVENT_TAKE_DAMAGE: dmg}
	if ok {
		l.Notify(takeDamageEvent)
	}
}
func (d *DummyGamestate) GetCardPicker() cards.AbstractCardPicker {
	return d.cardPiker
}
func (d *DummyGamestate) EndTurn() {
	// reset resource except money and reputation
	curRes := d.GetCurrentResource().Detail
	for k, _ := range curRes {
		if k == cards.RESOURCE_NAME_MONEY || k == cards.RESOURCE_NAME_REPUTATION {
			continue
		}
		d.GetCurrentResource().Detail[k] = 0
	}

	// remove cards played
	for _, c := range d.CardsPlayed {
		c.OnDiscarded()
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}
	d.CardsPlayed = []cards.Card{}

	// remove cards in hand
	for _, c := range d.CardsInHand {
		c.OnDiscarded()
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}
	d.CardsInHand = []cards.Card{}
	for _, c := range d.CenterCards {
		c.OnDiscarded()
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}
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
func (d *DummyGamestate) updateCenterCard(c cards.Card) {
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
func (d *DummyGamestate) Explore(c cards.Card) {
	// check cost and resource
	f := c.GetCost()
	res := d.currentResource
	if (&f).IsEnough(res) {
		// payResource
		d.PayResource(f)
		c.OnExplored()
		cardExploredEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_EXPLORED: c}

		l, ok := d.TopicsListeners[cards.EVENT_CARD_EXPLORED]
		if ok {
			l.Notify(cardExploredEvent)
		}
		// remove c from center cards
		d.updateCenterCard(c)
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
	f := c.GetCost()
	res := d.currentResource
	if (&f).IsEnough(res) {
		d.PayResource(f)
		c.OnSlain()

		cardDefeatedEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_DEFEATED: c}

		l, ok := d.TopicsListeners[cards.EVENT_CARD_DEFEATED]
		if ok {
			l.Notify(cardDefeatedEvent)
		}
		// remove c from center cards
		d.updateCenterCard(c)

	}
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
