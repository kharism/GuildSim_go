package cards_test

import (
	"fmt"
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

// show a basic card picker using stdin/stdout
type TestCardPicker struct {
	ChooseMethod func() int
}

func (t *TestCardPicker) PickCard(list []cards.Card, message string) int {
	fmt.Println(message)
	for i, card := range list {
		fmt.Printf("[%d] %s [%s] %s\n", i, card.GetName(), card.GetCost(), card.GetDescription())
	}
	// reader := bufio.NewReader(os.Stdin)
	// scanner := bufio.NewScanner(os.Stdin)

	// //text, _ := reader.ReadString('\n')
	// for scanner.Scan() {
	// 	picks, err := strconv.Atoi(scanner.Text())
	// 	if err != nil {
	// 		continue
	// 	}
	// 	return picks
	// }
	// return -1
	return t.ChooseMethod()
}

func NewRandomCardPickerChooser(list []cards.Card) func() int {
	return func() int {
		rand.Seed(10)
		return rand.Int() % len(list)
	}
}
func StaticCardPicker(index int) func() int {
	return func() int {
		return index
	}
}
func TestTextCardPicker(t *testing.T) {
	cardPicker := TestCardPicker{}
	cardPicker.ChooseMethod = StaticCardPicker(0)
	gamestate := NewDummyGamestate()
	rookieAdv := cards.NewRookieAdventurer(gamestate)
	packMule := cards.NewPackMule(gamestate)
	cards := []cards.Card{&rookieAdv, &packMule}
	cardsPicked := cardPicker.PickCard(cards, "Pick cards")
	t.Log(cardsPicked)
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
	CardsBanished     []cards.Card
	CardsDiscarded    cards.Deck
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
	d.cardPiker = &TestCardPicker{}
	d.CardsInDeck = cards.Deck{}
	d.CardsDiscarded = cards.Deck{}
	d.CardsBanished = []cards.Card{}

	d.HitPoint = 60
	return &d
}
func (d *DummyGamestate) GetCurrentHP() int {
	return d.HitPoint
}
func (d *DummyGamestate) TakeDamage(dmg int) {
	d.HitPoint -= dmg
	if dmg > 0 {
		l, ok := d.TopicsListeners[cards.EVENT_TAKE_DAMAGE]
		takeDamageEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_TAKE_DAMAGE_AMMOUNT: dmg}
		if ok {
			l.Notify(takeDamageEvent)
		}
	} else {
		l, ok := d.TopicsListeners[cards.EVENT_HEAL_DAMAGE]
		takeDamageEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_TAKE_DAMAGE_AMMOUNT: dmg}
		if ok {
			l.Notify(takeDamageEvent)
		}
	}

}
func (d *DummyGamestate) GetCardPicker() cards.AbstractCardPicker {
	return d.cardPiker
}
func (d *DummyGamestate) EndTurn() {
	// reset resource except money and reputation
	curRes := d.GetCurrentResource().Detail
	for k := range curRes {
		if k == cards.RESOURCE_NAME_MONEY || k == cards.RESOURCE_NAME_REPUTATION {
			continue
		}
		d.GetCurrentResource().Detail[k] = 0
	}

	// remove cards played
	for _, c := range d.CardsPlayed {
		// d.CardsDiscarded.Push(c)
		c.Dispose()
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}
	d.CardsPlayed = []cards.Card{}

	// remove cards in hand
	for _, c := range d.CardsInHand {
		// d.CardsDiscarded.Push(c)
		c.Dispose()
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}
	d.CardsInHand = []cards.Card{}
	for _, c := range d.CenterCards {
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}

}
func (d *DummyGamestate) AddCardToCenterDeck(c ...cards.Card) {
	for _, cc := range c {
		d.CardsInCenterDeck.Stack(cc)
	}
	d.CardsInCenterDeck.Shuffle()
}

// just play card from no particular location and added it to list of played card
// It will assume the card is played from hand and try to remove cards from hand if possible
// the card will not automatically go to discard/cooldown pile
// otherwise remove the card accordingly
func (d *DummyGamestate) PlayCard(c cards.Card) {
	d.RemoveCardFromHand(c)
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
	k := c.GetCost()
	if k.IsEnough(d.currentResource) {
		replacement := d.CardsInCenterDeck.Draw()
		d.RemoveCardFromCenterRow(c)
		d.CenterCards = append(d.CenterCards, replacement)
		if _, ok := c.(cards.Recruitable); ok {
			o := c.(cards.Recruitable)
			o.OnRecruit()
		}
		d.CardsDiscarded.Stack(c)
	}
	return
}
func (d *DummyGamestate) GetCooldownCard() []cards.Card {
	return d.CardsDiscarded.List()
}
func (d *DummyGamestate) DiscardCard(c cards.Card) {
	d.CardsDiscarded.Push(c)
	c.OnDiscarded()
	return
}
func (d *DummyGamestate) CenterRowInit() {
	f := d.ReplaceCenterCard()
	d.CenterCards = append(d.CenterCards, f)
}
func (d *DummyGamestate) RemoveCardFromHand(c cards.Card) {
	for idx, c2 := range d.CardsInHand {
		if c2 == c {
			d.RemoveCardFromHandIdx(idx)
			return
		}
	}
}
func (d *DummyGamestate) RemoveCardFromHandIdx(i int) {
	j := append(d.CardsInHand[:i], d.CardsInHand[i+1:]...)
	d.CardsInHand = j
}
func (d *DummyGamestate) RemoveCardFromCenterRow(c cards.Card) {
	for idx, c2 := range d.CenterCards {
		if c2 == c {
			d.RemoveCardFromCenterRowIdx(idx)
			return
		}
	}
}
func (d *DummyGamestate) RemoveCardFromCenterRowIdx(i int) {
	j := append(d.CenterCards[:i], d.CenterCards[i+1:]...)
	d.CenterCards = j
}
func (d *DummyGamestate) RemoveCardFromCooldown(c cards.Card) {
	for idx, c2 := range d.CardsDiscarded.List() {
		if c2 == c {
			d.RemoveCardFromCooldownIdx(idx)
			return
		}
	}
}
func (d *DummyGamestate) RemoveCardFromCooldownIdx(i int) {
	cooldownList := d.CardsDiscarded.List()
	j := append(cooldownList[:i], cooldownList[i+1:]...)
	d.CardsDiscarded.SetList(j)
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
	if d.CardsInDeck.Size() == 0 {
		// shuffle discard pile
		d.CardsDiscarded.Shuffle()
		d.CardsInDeck = d.CardsDiscarded
		d.CardsDiscarded = cards.Deck{}
	}
	newCard := d.CardsInDeck.Draw()
	d.CardsInHand = append(d.CardsInHand, newCard)
	newCard.OnAddedToHand()
	return
}
func (d *DummyGamestate) BanishCard(c cards.Card) {
	d.CardsBanished = append(d.CardsBanished, c)
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_BANISHED]; ok {
		notification := map[string]interface{}{cards.EVENT_ATTR_CARD_BANISHED: c}
		d.TopicsListeners[cards.EVENT_CARD_BANISHED].Notify(notification)
	}
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
