package gamestate

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/cards/item"
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
	RuleEnforcer      map[string]*cards.RuleEnforcer
	CardsInHand       []cards.Card
	CardsPlayed       []cards.Card
	CenterCards       []cards.Card
	HitPoint          int
	CardsDiscarded    cards.Deck
	CardsBanished     []cards.Card
	ItemCards         []cards.Card
	//ui stuff
	cardPiker         cards.AbstractCardPicker
	centerCardChanged bool
	boolPicker        cards.AbstractBoolPicker
	cardViewer        cards.AbstractDetailViewer
}

// AddCardToCenterDeck implements cards.AbstractGamestate
func (d *DefaultGamestate) AddCardToCenterDeck(source string, shuffle bool, c ...cards.Card) {
	var j *DummyEventListener
	j = d.TopicsListeners[cards.EVENT_CARD_GOTO_CENTER]
	for _, cc := range c {
		d.CardsInCenterDeck.Stack(cc)
		if j != nil {
			evt := map[string]interface{}{cards.EVENT_ATTR_CARD_GOTO_CENTER: cc, cards.EVENT_ATTR_DISCARD_SOURCE: source}
			j.Notify(evt)
		}
	}
	if shuffle {
		d.CardsInCenterDeck.Shuffle()
	}

}

func NewDefaultGamestate() cards.AbstractGamestate {
	d := DefaultGamestate{}
	d.currentResource = cards.NewResource()
	d.CardsPlayed = []cards.Card{}
	d.TopicsListeners = map[string]*DummyEventListener{}
	d.RuleEnforcer = map[string]*cards.RuleEnforcer{}
	d.CenterCards = []cards.Card{}
	d.CardsInHand = []cards.Card{}
	d.CardsBanished = []cards.Card{}
	d.ItemCards = []cards.Card{}
	// d.cardPiker = &TextCardPicker{}
	d.HitPoint = 60
	d.CardsDiscarded = cards.Deck{}
	d.CardsInCenterDeck = cards.Deck{}
	d.CardsInDeck = cards.Deck{}
	return &d
}
func (d *DefaultGamestate) AttachLegalCheck(actionName string, lc cards.LegalChecker) {
	if _, ok := d.RuleEnforcer[actionName]; !ok {
		d.RuleEnforcer[actionName] = cards.NewRuleEnforcer()
	}
	d.RuleEnforcer[actionName].AttachRule(lc)
}
func (d *DefaultGamestate) DetachLegalCheck(actionName string, lc cards.LegalChecker) {
	if _, ok := d.RuleEnforcer[actionName]; !ok {
		return
	}
	d.RuleEnforcer[actionName].DetachRule(lc)
}
func (d *DefaultGamestate) LegalCheck(actionName string, data interface{}) bool {
	j, ok := d.RuleEnforcer[actionName]
	if !ok {
		return true
	}
	return j.Check(data)
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
func (d *DefaultGamestate) ListItems() []cards.Card {
	return d.ItemCards
}
func (d *DefaultGamestate) RemoveItem(c cards.Card) {
	idx := -1
	for i := range d.ItemCards {
		if d.ItemCards[i] == c {
			idx = i
			break
		}
	}
	if idx != -1 {
		d.RemoveItemIndex(idx)
		data := map[string]interface{}{cards.EVENT_ATTR_ITEM_REMOVED: c}
		d.NotifyListener(cards.EVENT_ITEM_REMOVED, data)
	}
}
func (d *DefaultGamestate) GenerateRandomPotion(rarity int) cards.Card {
	return item.CreatePotionRandom(d, rarity)
}
func (d *DefaultGamestate) GenerateRandomRelic(rarity int) cards.Card {
	return item.CreateRelicRandom(d, rarity)
}
func (d *DefaultGamestate) RemoveItemIndex(i int) {
	h := append(d.ItemCards[:i], d.ItemCards[i+1:]...)
	d.ItemCards = h
}
func (d *DefaultGamestate) AddItem(c cards.Card) {
	fmt.Print(c.GetName())
	d.ItemCards = append(d.ItemCards, c)
	c.OnAcquire()
	data := map[string]interface{}{cards.EVENT_ATTR_ITEM_ADDED: c}
	d.NotifyListener(cards.EVENT_ITEM_ADDED, data)
}
func (d *DefaultGamestate) NotifyListener(eventname string, data map[string]interface{}) {
	if _, ok := d.TopicsListeners[eventname]; ok {
		j := d.TopicsListeners[eventname]
		j.Notify(data)
	}
}
func (d *DefaultGamestate) ConsumeItem(c cards.Consumable) {
	c.OnConsume()
}
func (d *DefaultGamestate) GetCurrentHP() int {
	return d.HitPoint
}
func (d *DefaultGamestate) TakeDamage(dmg int) {
	block, ok := d.GetCurrentResource().Detail[cards.RESOURCE_NAME_BLOCK]
	if ok {
		if dmg <= block {
			// d.PayResource()
			d.GetCurrentResource().Detail[cards.RESOURCE_NAME_BLOCK] -= dmg
			dmg = 0
		} else {
			dmg -= block
			d.GetCurrentResource().Detail[cards.RESOURCE_NAME_BLOCK] = 0
		}

	}

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

func (d *DefaultGamestate) GetCardPicker() cards.AbstractCardPicker {
	return d.cardPiker
}
func (d *DefaultGamestate) SetCardPicker(a cards.AbstractCardPicker) {
	d.cardPiker = a
}
func (d *DefaultGamestate) GetBoolPicker() cards.AbstractBoolPicker {
	return d.boolPicker
}
func (d *DefaultGamestate) SetBoolPicker(a cards.AbstractBoolPicker) {
	d.boolPicker = a
}
func (d *DefaultGamestate) EndTurn() {

	// remove cards played
	for _, c := range d.CardsPlayed {
		c.Dispose(cards.DISCARD_SOURCE_PLAYED)
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}
	d.CardsPlayed = []cards.Card{}

	// remove cards in hand
	for _, c := range d.CardsInHand {
		c.Dispose(cards.DISCARD_SOURCE_HAND)
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}
	// reset resource except money and reputation
	curRes := d.GetCurrentResource().Detail
	for k := range curRes {
		if k == cards.RESOURCE_NAME_MONEY || k == cards.RESOURCE_NAME_REPUTATION {
			continue
		}
		d.GetCurrentResource().Detail[k] = 0
	}
	if !d.centerCardChanged {
		cardsShuffledBack := 0
		for i := len(d.CenterCards) - 1; i >= 0; i-- {
			hh := d.CenterCards[i]
			if _, ok := hh.(cards.Unshuffleable); ok {
				continue
			}
			cardsShuffledBack++
			d.RemoveCardFromCenterRowIdx(i)
			d.AddCardToCenterDeck(cards.DISCARD_SOURCE_CENTER, false, hh)
		}
		d.CardsInCenterDeck.Shuffle()
		for i := 0; i < cardsShuffledBack; i++ {
			f := d.ReplaceCenterCard()
			d.CenterCards = append(d.CenterCards, f)
		}
	}
	d.CardsInHand = []cards.Card{}
	for _, c := range d.CenterCards {
		// c.OnDiscarded()
		fmt.Println("Check Punish", c.GetName())
		if pun, ok := c.(cards.Punisher); ok {
			fmt.Println("Punish")
			pun.OnPunish()
		}
	}
	if _, ok := d.TopicsListeners[cards.EVENT_END_OF_TURN]; ok {
		j := d.TopicsListeners[cards.EVENT_END_OF_TURN]
		data := map[string]interface{}{}
		j.Notify(data)
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
	if !d.LegalCheck(cards.ACTION_RECRUIT, c) {
		return
	}
	k := c.GetCost()
	if k.IsEnough(d.currentResource) {
		d.PayResource(k)
		if _, ok := c.(cards.Recruitable); ok {
			o := c.(cards.Recruitable)
			o.OnRecruit()
		}
		fmt.Println("DefGS Recruit")
		d.RemoveCardFromCenterRow(c)
		d.CardsDiscarded.Stack(c)
		d.updateCenterCard(c)
		if _, ok := d.TopicsListeners[cards.EVENT_CARD_RECRUITED]; ok {
			evtDetails := map[string]interface{}{cards.EVENT_ATTR_CARD_RECRUITED: c}
			j := d.TopicsListeners[cards.EVENT_CARD_RECRUITED]
			j.Notify(evtDetails)
		}

		// d.CenterCards = append(d.CenterCards, replacement)

	}
	return
}
func (d *DefaultGamestate) GetCooldownCard() []cards.Card {
	return d.CardsDiscarded.List()
}
func (d *DefaultGamestate) DiscardCard(c cards.Card, source string) {
	d.CardsDiscarded.Push(c)
	c.OnDiscarded()
	data := map[string]interface{}{cards.EVENT_ATTR_CARD_DISCARDED: c, cards.EVENT_ATTR_DISCARD_SOURCE: source}
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_DISCARDED]; ok {
		d.TopicsListeners[cards.EVENT_CARD_DISCARDED].Notify(data)
	}
	return
}
func (d *DefaultGamestate) CenterRowInit() {
	d.CardsInCenterDeck.Shuffle()
	for i := 0; i < 5; i++ {
		f := d.ReplaceCenterCard()
		d.CenterCards = append(d.CenterCards, f)
	}
}
func (d *DefaultGamestate) StackCards(source string, cc ...cards.Card) {
	for _, c := range cc {
		d.CardsInDeck.Stack(c)
		if _, ok := d.TopicsListeners[cards.EVENT_ATTR_CARD_STACKED]; ok {
			j := d.TopicsListeners[cards.EVENT_ATTR_CARD_STACKED]
			data := map[string]interface{}{cards.EVENT_ATTR_CARD_STACKED: c, cards.EVENT_ATTR_DISCARD_SOURCE: source}
			j.Notify(data)
		}
	}

}
func (d *DefaultGamestate) ShuffleMainDeck() {
	d.CardsInDeck.Shuffle()
}
func (d *DefaultGamestate) updateCenterCard(c cards.Card) {
	replacementCard := d.ReplaceCenterCard()
	newCenterCards := []cards.Card{}
	fmt.Println("Replace", c.GetName(), "with", replacementCard.GetName())
	for _, v := range d.CenterCards {
		if v == c {
			newCenterCards = append(newCenterCards, replacementCard)
		} else {
			newCenterCards = append(newCenterCards, v)
		}
	}
	newCenterCards = append(newCenterCards, replacementCard)
	d.CenterCards = newCenterCards

}
func (d *DefaultGamestate) Explore(c cards.Card) {
	// check cost and resource
	f := c.GetCost()
	res := d.currentResource
	if (&f).IsEnough(res) {
		// payResource
		d.PayResource(f)
		c.OnExplored()
		d.RemoveCardFromCenterRow(c)
		// remove c from center cards
		d.updateCenterCard(c)
		d.BanishCard(c, cards.DISCARD_SOURCE_CENTER)
		cardExploredEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_EXPLORED: c}

		l, ok := d.TopicsListeners[cards.EVENT_CARD_EXPLORED]
		if ok {
			l.Notify(cardExploredEvent)
		}

	}
}
func (d *DefaultGamestate) AppendCenterCard(c cards.Card) {
	d.CenterCards = append(d.CenterCards, c)
}
func (d *DefaultGamestate) SetDetailViewer(v cards.AbstractDetailViewer) {
	d.cardViewer = v
}
func (d *DefaultGamestate) GetDetailViewer() cards.AbstractDetailViewer {
	return d.cardViewer
}
func (d *DefaultGamestate) PeekCenterCard() cards.Card {
	if len(d.CardsInCenterDeck.List()) == 0 {
		return nil
	}
	return d.CardsInCenterDeck.List()[0]
}
func (d *DefaultGamestate) ReplaceCenterCard() cards.Card {
	d.centerCardChanged = true
	replacementCard := d.CardsInCenterDeck.Draw()
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_DRAWN_CENTER]; ok {
		evtDetails := map[string]interface{}{cards.EVENT_ATTR_CARD_DRAWN: replacementCard}
		j := d.TopicsListeners[cards.EVENT_CARD_DRAWN_CENTER]
		j.Notify(evtDetails)
	}
	if _, ok := replacementCard.(cards.Trapper); ok {
		j := replacementCard.(cards.Trapper)
		if !j.IsDisarmed() {
			j.Trap()
		} else {
			d.BanishCard(replacementCard, cards.DISCARD_SOURCE_CENTER)
			replacementCard = d.CardsInCenterDeck.Draw()
			evtDetails := map[string]interface{}{cards.EVENT_ATTR_CARD_DRAWN: replacementCard}
			d.NotifyListener(cards.EVENT_CARD_DRAWN_CENTER, evtDetails)
		}

	}
	return replacementCard
}
func (d *DefaultGamestate) BeginTurn() {
	d.centerCardChanged = false
	for i := 0; i < 5; i++ {
		d.Draw()
	}
	if _, ok := d.TopicsListeners[cards.EVENT_START_OF_TURN]; ok {
		j := d.TopicsListeners[cards.EVENT_START_OF_TURN]
		data := map[string]interface{}{}
		j.Notify(data)
	}
}
func (d *DefaultGamestate) Draw() {
	if !d.LegalCheck(cards.ACTION_DRAW, nil) {
		return
	}
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
func (d *DefaultGamestate) BanishCard(c cards.Card, source string) {
	d.CardsBanished = append(d.CardsBanished, c)
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_BANISHED]; ok {
		notification := map[string]interface{}{cards.EVENT_ATTR_CARD_BANISHED: c, cards.EVENT_ATTR_DISCARD_SOURCE: source}
		d.TopicsListeners[cards.EVENT_CARD_BANISHED].Notify(notification)
	}
	return
}
func (d *DefaultGamestate) Disarm(c cards.Card) {
	f := c.GetCost()
	res := d.currentResource
	if (&f).IsEnough(res) {
		d.PayResource(f)
		c.OnSlain()
		d.RemoveCardFromCenterRow(c)
		// remove c from center cards
		d.updateCenterCard(c)
		// d.BanishCard(c, cards.DISCARD_SOURCE_CENTER)
		c.Dispose(cards.DISCARD_SOURCE_CENTER)
		trapRemovedEvent := map[string]interface{}{cards.EVENT_ATTR_TRAP_REMOVED: c}
		d.NotifyListener(cards.EVENT_TRAP_REMOVED, trapRemovedEvent)

	}
	return
}
func (d *DefaultGamestate) DefeatCard(c cards.Card) {
	if !d.LegalCheck(cards.ACTION_DEFEAT, c) {
		return
	}
	f := c.GetCost()
	res := d.currentResource
	if (&f).IsEnough(res) {
		d.PayResource(f)
		c.OnSlain()
		d.RemoveCardFromCenterRow(c)
		// remove c from center cards
		d.updateCenterCard(c)
		d.BanishCard(c, cards.DISCARD_SOURCE_CENTER)
		cardDefeatedEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_DEFEATED: c}

		l, ok := d.TopicsListeners[cards.EVENT_CARD_DEFEATED]
		if ok {
			l.Notify(cardDefeatedEvent)
		}

	}
	return
}
func (d *DefaultGamestate) GetCurrentResource() cards.Resource {
	return d.currentResource
}
func (d *DefaultGamestate) AddResource(name string, amount int) {
	d.currentResource.AddResource(name, amount)
}
