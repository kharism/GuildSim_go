package gamestate

import (
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/observer"
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

type DefaultGamestate struct {
	currentResource   cards.Resource
	CardsInDeck       cards.Deck
	CardsInCenterDeck cards.Deck
	TopicsListeners   map[string]*DummyEventListener
	CardsInHand       []cards.Card
	CardsPlayed       []cards.Card
	CenterCards       []cards.Card
	HitPoint          int
	CardsDiscarded    cards.Deck
	CardsBanished     []cards.Card
	//ui stuff
	cardPiker cards.AbstractCardPicker
}

// AddCardToCenterDeck implements cards.AbstractGamestate
func (d *DefaultGamestate) AddCardToCenterDeck(c ...cards.Card) {
	for _, cc := range c {
		d.CardsInCenterDeck.Stack(cc)
	}
	d.CardsInCenterDeck.Shuffle()
}

func NewDefaultGamestate() cards.AbstractGamestate {
	d := DefaultGamestate{}
	d.currentResource = cards.NewResource()
	d.CardsPlayed = []cards.Card{}
	d.TopicsListeners = map[string]*DummyEventListener{}
	d.CenterCards = []cards.Card{}
	d.CardsInHand = []cards.Card{}
	d.CardsBanished = []cards.Card{}
	// d.cardPiker = &TextCardPicker{}
	d.HitPoint = 60
	d.CardsDiscarded = cards.Deck{}
	d.CardsInCenterDeck = cards.Deck{}
	d.CardsInDeck = cards.Deck{}
	return &d
}

func (d *DefaultGamestate) PayResource(cost cards.Cost) {
	for key, val := range cost.Detail {
		d.currentResource.Detail[key] -= val
	}
}
func (d *DefaultGamestate) RemoveCardFromHand(c cards.Card) {
	for idx, c2 := range d.CardsInHand {
		if c2 == c {
			d.RemoveCardFromHandIdx(idx)
			return
		}
	}
}
func (d *DefaultGamestate) RemoveCardFromHandIdx(i int) {
	j := append(d.CardsInHand[:i], d.CardsInHand[i+1:]...)
	d.CardsInHand = j
}
func (d *DefaultGamestate) RemoveCardFromCenterRow(c cards.Card) {
	for idx, c2 := range d.CenterCards {
		if c2 == c {
			d.RemoveCardFromCenterRowIdx(idx)
			return
		}
	}
}
func (d *DefaultGamestate) RemoveCardFromCenterRowIdx(i int) {
	j := append(d.CenterCards[:i], d.CenterCards[i+1:]...)
	d.CenterCards = j
}
func (d *DefaultGamestate) RemoveCardFromCooldown(c cards.Card) {
	for idx, c2 := range d.CardsDiscarded.List() {
		if c2 == c {
			d.RemoveCardFromCooldownIdx(idx)
			return
		}
	}
}
func (d *DefaultGamestate) RemoveCardFromCooldownIdx(i int) {
	cooldownList := d.CardsDiscarded.List()
	j := append(cooldownList[:i], cooldownList[i+1:]...)
	d.CardsDiscarded.SetList(j)
}
func (d *DefaultGamestate) AttachListener(eventName string, l observer.Listener) {
	if _, ok := d.TopicsListeners[eventName]; !ok {
		d.TopicsListeners[eventName] = &DummyEventListener{}
	}
	k := (d.TopicsListeners[eventName])
	k.Attach(l)
	// fmt.Println("Attach Listener", len(d.TopicsListeners[eventName].Listeners))
}
func (d *DefaultGamestate) RemoveListener(eventName string, l observer.Listener) {
	if _, ok := d.TopicsListeners[eventName]; !ok {
		return
	}
	// fmt.Println("Remove Listener")
	k := (d.TopicsListeners[eventName])
	k.Detach(l)
}

func (d *DefaultGamestate) GetCurrentHP() int {
	return d.HitPoint
}
func (d *DefaultGamestate) TakeDamage(dmg int) {
	d.HitPoint -= dmg
	l, ok := d.TopicsListeners[cards.EVENT_TAKE_DAMAGE]
	takeDamageEvent := map[string]interface{}{cards.EVENT_TAKE_DAMAGE: dmg}
	if ok {
		l.Notify(takeDamageEvent)
	}
}

func (d *DefaultGamestate) GetCardPicker() cards.AbstractCardPicker {
	return d.cardPiker
}
func (d *DefaultGamestate) SetCardPicker(a cards.AbstractCardPicker) {
	d.cardPiker = a
}
func (d *DefaultGamestate) EndTurn() {
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
		c.Dispose()
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}
	d.CardsPlayed = []cards.Card{}

	// remove cards in hand
	for _, c := range d.CardsInHand {
		c.Dispose()
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

func (d *DefaultGamestate) PlayCard(c cards.Card) {
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
func (d *DefaultGamestate) GetPlayedCards() []cards.Card {
	return d.CardsPlayed
}
func (d *DefaultGamestate) GetCardInHand() []cards.Card {
	return d.CardsInHand
}
func (d *DefaultGamestate) GetCenterCard() []cards.Card {
	return d.CenterCards
}
func (d *DefaultGamestate) RecruitCard(c cards.Card) {
	k := c.GetCost()
	if k.IsEnough(d.currentResource) {
		d.PayResource(k)
		replacement := d.CardsInCenterDeck.Draw()
		d.RemoveCardFromCenterRow(c)
		d.CenterCards = append(d.CenterCards, replacement)
		if _, ok := c.(cards.Recruitable); ok {
			o := c.(cards.Recruitable)
			o.OnRecruit()
		}
		if _, ok := d.TopicsListeners[cards.EVENT_CARD_RECRUITED]; ok {
			evtDetails := map[string]interface{}{cards.EVENT_ATTR_CARD_RECRUITED: c}
			j := d.TopicsListeners[cards.EVENT_CARD_RECRUITED]
			j.Notify(evtDetails)
		}
		d.CardsDiscarded.Stack(c)
	}
	return
}
func (d *DefaultGamestate) GetCooldownCard() []cards.Card {
	return d.CardsDiscarded.List()
}
func (d *DefaultGamestate) DiscardCard(c cards.Card) {
	d.CardsDiscarded.Push(c)
	c.OnDiscarded()
	return
}
func (d *DefaultGamestate) CenterRowInit() {
	d.CardsInCenterDeck.Shuffle()
	for i := 0; i < 5; i++ {
		f := d.ReplaceCenterCard()
		d.CenterCards = append(d.CenterCards, f)
		if _, ok := d.TopicsListeners[cards.EVENT_CARD_DRAWN_CENTER]; ok {
			evtDetails := map[string]interface{}{cards.EVENT_ATTR_CARD_DRAWN: f}
			j := d.TopicsListeners[cards.EVENT_CARD_DRAWN_CENTER]
			j.Notify(evtDetails)
		}
	}
}
func (d *DefaultGamestate) updateCenterCard(c cards.Card) {
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
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_DRAWN_CENTER]; ok {
		evtDetails := map[string]interface{}{cards.EVENT_ATTR_CARD_DRAWN: replacementCard}
		j := d.TopicsListeners[cards.EVENT_CARD_DRAWN_CENTER]
		j.Notify(evtDetails)
	}
}
func (d *DefaultGamestate) Explore(c cards.Card) {
	// check cost and resource
	f := c.GetCost()
	res := d.currentResource
	if (&f).IsEnough(res) {
		// payResource
		d.PayResource(f)
		c.OnExplored()
		d.BanishCard(c)
		cardExploredEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_EXPLORED: c}

		l, ok := d.TopicsListeners[cards.EVENT_CARD_EXPLORED]
		if ok {
			l.Notify(cardExploredEvent)
		}
		// remove c from center cards
		d.updateCenterCard(c)
	}
}
func (d *DefaultGamestate) ReplaceCenterCard() cards.Card {
	return d.CardsInCenterDeck.Draw()
}
func (d *DefaultGamestate) BeginTurn() {
	for i := 0; i < 5; i++ {
		d.Draw()
	}
}
func (d *DefaultGamestate) Draw() {
	if d.CardsInDeck.Size() == 0 {
		// shuffle discard pile
		d.CardsDiscarded.Shuffle()
		d.CardsInDeck = d.CardsDiscarded
		d.CardsDiscarded = cards.Deck{}
	}
	newCard := d.CardsInDeck.Draw()
	d.CardsInHand = append(d.CardsInHand, newCard)
	newCard.OnAddedToHand()
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_DRAWN]; ok {
		evtDetails := map[string]interface{}{cards.EVENT_ATTR_CARD_DRAWN: newCard}
		j := d.TopicsListeners[cards.EVENT_CARD_DRAWN]
		j.Notify(evtDetails)
	}
	return
}
func (d *DefaultGamestate) BanishCard(c cards.Card) {
	d.CardsBanished = append(d.CardsBanished, c)
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_BANISHED]; ok {
		notification := map[string]interface{}{cards.EVENT_ATTR_CARD_BANISHED: c}
		d.TopicsListeners[cards.EVENT_CARD_BANISHED].Notify(notification)
	}
	return
}
func (d *DefaultGamestate) DefeatCard(c cards.Card) {
	f := c.GetCost()
	res := d.currentResource
	if (&f).IsEnough(res) {
		d.PayResource(f)
		c.OnSlain()
		d.BanishCard(c)
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
func (d *DefaultGamestate) GetCurrentResource() cards.Resource {
	return d.currentResource
}
func (d *DefaultGamestate) AddResource(name string, amount int) {
	d.currentResource.AddResource(name, amount)
}
