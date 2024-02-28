package gamestate

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/cards/item"
	"github/kharism/GuildSim_go/internal/factory"
	"github/kharism/GuildSim_go/internal/observer"
	"sync"
)

type DummyEventListener struct {
	mutex     *sync.Mutex
	Listeners []observer.Listener
}

func NewDummyEventListener() *DummyEventListener {
	d := DummyEventListener{}
	d.mutex = &sync.Mutex{}
	d.Listeners = []observer.Listener{}
	return &d
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
	d.mutex.Lock()
	for _, i := range d.Listeners {
		i.DoAction(data)
	}
	d.mutex.Unlock()
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

	CurrentFillerIdx int
	FillerFuncList   []func(cards.AbstractGamestate) []cards.Card

	//ui stuff
	cardPiker         cards.AbstractCardPicker
	centerCardChanged bool
	boolPicker        cards.AbstractBoolPicker
	cardViewer        cards.AbstractDetailViewer
	mutex             *sync.Mutex

	// act manager
	ActDecorator []func(cards.AbstractGamestate) cards.AbstractGamestate
}

func (d *DefaultGamestate) MutexLock() {
	d.mutex.Lock()
}
func (d *DefaultGamestate) MutexUnlock() {
	d.mutex.Unlock()
}

// AddCardToCenterDeck implements cards.AbstractGamestate
func (d *DefaultGamestate) AddCardToCenterDeck(source string, shuffle bool, c ...cards.Card) {
	j := d.TopicsListeners[cards.EVENT_CARD_GOTO_CENTER]
	for _, cc := range c {
		d.mutex.Lock()
		d.CardsInCenterDeck.Stack(cc)
		d.mutex.Unlock()
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
	d.mutex = &sync.Mutex{}
	d.FillerFuncList = append(d.FillerFuncList, factory.CreateFillerCenterDeck1, factory.CreateFillerCenterDeck2)
	return &d
}
func (d *DefaultGamestate) SetFillerIndex(i int) {
	d.CurrentFillerIdx = i
}
func (d *DefaultGamestate) GetFillerIndex() int {
	return d.CurrentFillerIdx
}
func (d *DefaultGamestate) AppendCardFiller(newFunc func(cards.AbstractGamestate) []cards.Card) {
	d.FillerFuncList = append(d.FillerFuncList, newFunc)
}
func (d *DefaultGamestate) AttachLegalCheck(actionName string, lc cards.LegalChecker) {
	if _, ok := d.RuleEnforcer[actionName]; !ok {
		d.RuleEnforcer[actionName] = cards.NewRuleEnforcer()
	}
	d.RuleEnforcer[actionName].AttachRule(lc)
	data := map[string]interface{}{cards.EVENT_ATTR_LIMITER: lc, cards.EVENT_ATTR_LIMITER_ACTION: actionName}
	fmt.Println("attach RuleCheck")
	d.NotifyListener(cards.EVENT_ATTACH_LIMITER, data)
}
func (d *DefaultGamestate) DetachLegalCheck(actionName string, lc cards.LegalChecker) {
	if _, ok := d.RuleEnforcer[actionName]; !ok {
		return
	}
	d.RuleEnforcer[actionName].DetachRule(lc)
	data := map[string]interface{}{cards.EVENT_ATTR_LIMITER: lc, cards.EVENT_ATTR_LIMITER_ACTION: actionName}
	d.NotifyListener(cards.EVENT_DETACH_LIMITER, data)
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
		//d.currentResource.Detail[key] -= val
		d.AddResource(key, -val)
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
func (d *DefaultGamestate) RemoveCardFromPlayed(c cards.Card) {
	for idx, c2 := range d.CardsPlayed {
		if c2 == c {
			d.RemoveCardFromHandIdx(idx)
			return
		}
	}

}
func (d *DefaultGamestate) RemoveCardFromPlayedIdx(i int) {
	j := append(d.CardsPlayed[:i], d.CardsPlayed[i+1:]...)
	d.CardsPlayed = j
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
		d.TopicsListeners[eventName] = NewDummyEventListener()
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
func (d *DefaultGamestate) ActDecorators() []func(cards.AbstractGamestate) cards.AbstractGamestate {
	return d.ActDecorator
}
func (d *DefaultGamestate) AddActDecorator(f func(cards.AbstractGamestate) cards.AbstractGamestate) {
	d.ActDecorator = append(d.ActDecorator, f)
}
func (d *DefaultGamestate) ConsumeItem(c cards.Consumable) {
	c.OnConsume()
}
func (d *DefaultGamestate) GetCurrentHP() int {
	return d.HitPoint
}
func (d *DefaultGamestate) TakeDamage(dmg int) {
	if dmg > 0 {
		block, ok := d.GetCurrentResource().Detail[cards.RESOURCE_NAME_BLOCK]
		if ok {
			if dmg <= block {
				// d.PayResource()
				// d.GetCurrentResource().Detail[cards.RESOURCE_NAME_BLOCK] -= dmg
				cost := cards.NewCost()
				cost.AddResource(cards.RESOURCE_NAME_BLOCK, dmg)
				d.PayResource(cost)
				dmg = 0
			} else {
				dmg -= block
				d.GetCurrentResource().Detail[cards.RESOURCE_NAME_BLOCK] = 0
			}

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

	// remove cards in hand
	for i := len(d.CardsInHand) - 1; i >= 0; i-- {
		c := d.CardsInHand[i]
		c.Dispose(cards.DISCARD_SOURCE_HAND)
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
		d.RemoveCardFromHand(c)
	}
	// reset resource except money and reputation
	d.MutexLock()
	curRes := d.GetCurrentResource().Detail
	for k := range curRes {
		if k == cards.RESOURCE_NAME_MONEY || k == cards.RESOURCE_NAME_REPUTATION || k == cards.RESOURCE_NAME_BLOCK {
			continue
		}
		decreaseAmount := d.GetCurrentResource().Detail[k]
		d.AddResource(k, -decreaseAmount)
		//d.GetCurrentResource().Detail[k] -= decreaseAmount

	}
	d.MutexUnlock()
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
			fmt.Println("Replace center card with", f.GetName())
			d.CenterCards = append(d.CenterCards, f)
		}
	}
	d.CardsInHand = []cards.Card{}
	for _, c := range d.CenterCards {
		// c.OnDiscarded()
		fmt.Println("Check Punish", c.GetName())
		if pun, ok := c.(cards.Punisher); ok {
			fmt.Println("Punish")
			// TODO: add pre-punish animation
			if _, ok := d.TopicsListeners[cards.EVENT_BEFORE_PUNISH]; ok {
				j := d.TopicsListeners[cards.EVENT_BEFORE_PUNISH]
				data := map[string]interface{}{}
				data[cards.EVENT_ATTR_BEFORE_PUNISH_CARD] = c
				j.Notify(data)
			}
			pun.OnPunish()
		}
	}
	d.mutex.Lock()
	// remove cards played
	for i := len(d.CardsPlayed) - 1; i >= 0; i-- {
		c := d.CardsPlayed[i]
		c.Dispose(cards.DISCARD_SOURCE_PLAYED)
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}
	d.CardsPlayed = []cards.Card{}
	d.mutex.Unlock()
	// d.GetCurrentResource().Detail[cards.RESOURCE_NAME_BLOCK] = 0
	decreaseAmount := d.GetCurrentResource().Detail[cards.RESOURCE_NAME_BLOCK]
	d.AddResource(cards.RESOURCE_NAME_BLOCK, -decreaseAmount)
	if _, ok := d.TopicsListeners[cards.EVENT_END_OF_TURN]; ok {
		j := d.TopicsListeners[cards.EVENT_END_OF_TURN]
		data := map[string]interface{}{}
		j.Notify(data)
	}

}

func (d *DefaultGamestate) PlayCard(c cards.Card) {
	d.mutex.Lock()
	d.RemoveCardFromHand(c)
	d.mutex.Unlock()
	c.OnPlay()
	// fmt.Println("Card played", c.GetName())
	cardPlayedEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_PLAYED: c}

	l, ok := d.TopicsListeners[cards.EVENT_CARD_PLAYED]
	if ok {
		l.Notify(cardPlayedEvent)
	}
	d.mutex.Lock()
	d.CardsPlayed = append(d.CardsPlayed, c)
	d.mutex.Unlock()
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
	d.mutex.Lock()
	isEnough := k.IsEnough(d.currentResource)
	d.mutex.Unlock()
	if isEnough {
		d.mutex.Lock()
		d.PayResource(k)
		d.mutex.Unlock()
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
}
func (d *DefaultGamestate) CenterRowInit() {
	// d.CardsInCenterDeck.Shuffle()
	for i := 0; i < 5; i++ {
		f := d.ReplaceCenterCard()
		d.CenterCards = append(d.CenterCards, f)
	}
}
func (d *DefaultGamestate) StackCards(source string, cc ...cards.Card) {
	for _, c := range cc {
		d.CardsInDeck.Stack(c)
		c.OnReturnToDeck()
		if _, ok := d.TopicsListeners[cards.EVENT_CARD_STACKED]; ok {
			j := d.TopicsListeners[cards.EVENT_CARD_STACKED]
			data := map[string]interface{}{cards.EVENT_ATTR_CARD_STACKED: c, cards.EVENT_ATTR_DISCARD_SOURCE: source}
			j.Notify(data)
		}
	}

}
func (d *DefaultGamestate) GetMainDeck() *cards.Deck {
	return &d.CardsInDeck
}
func (d *DefaultGamestate) ShuffleMainDeck() {
	d.CardsInDeck.Shuffle()
}
func (d *DefaultGamestate) UpdateCenterCard(c cards.Card) {
	d.updateCenterCard(c)
}
func (d *DefaultGamestate) updateCenterCard(c cards.Card) {
	replacementCard := d.ReplaceCenterCard()
	if d.CardsInCenterDeck.Size() == 0 {
		filler := d.FillerFuncList[d.CurrentFillerIdx](d) //factory.CardFactory(factory.SET_FILLER_CARDS, d)
		d.CardsInCenterDeck.SetList(filler)
	}
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
	d.mutex.Lock()
	isEnough := (&f).IsEnough(res)
	d.mutex.Unlock()
	if isEnough {
		// payResource
		d.mutex.Lock()
		d.PayResource(f)
		d.mutex.Unlock()
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
	d.mutex.Lock()
	replacementCard := d.CardsInCenterDeck.Draw()
	d.mutex.Unlock()
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_DRAWN_CENTER]; ok {
		evtDetails := map[string]interface{}{cards.EVENT_ATTR_CARD_DRAWN: replacementCard}
		j := d.TopicsListeners[cards.EVENT_CARD_DRAWN_CENTER]
		j.Notify(evtDetails)
	}
	if _, ok := replacementCard.(cards.Trapper); ok {
		j := replacementCard.(cards.Trapper)
		if !j.IsDisarmed() {
			fmt.Println("Invoke Trap")
			data := map[string]interface{}{}
			data[cards.EVENT_ATTR_BEFORE_TRAP] = j
			d.NotifyListener(cards.EVENT_BEFORE_TRAP, data)
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
		fmt.Println("Illegal to draw")
		return
	}
	if d.CardsInDeck.Size() == 0 {
		// shuffle discard pile
		d.CardsDiscarded.Shuffle()
		d.CardsInDeck = d.CardsDiscarded
		d.CardsDiscarded = cards.Deck{}
		// if still no cards in deck after joining discard and deck, return
		if d.CardsInDeck.Size() == 0 {
			return
		}
	}
	newCard := d.CardsInDeck.Draw()
	d.CardsInHand = append(d.CardsInHand, newCard)
	newCard.OnAddedToHand()
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_DRAWN]; ok {
		evtDetails := map[string]interface{}{cards.EVENT_ATTR_CARD_DRAWN: newCard}
		j := d.TopicsListeners[cards.EVENT_CARD_DRAWN]
		j.Notify(evtDetails)
	}

}
func (d *DefaultGamestate) BanishCard(c cards.Card, source string) {
	d.CardsBanished = append(d.CardsBanished, c)
	c.OnBanished()
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_BANISHED]; ok {
		notification := map[string]interface{}{cards.EVENT_ATTR_CARD_BANISHED: c, cards.EVENT_ATTR_DISCARD_SOURCE: source}
		d.TopicsListeners[cards.EVENT_CARD_BANISHED].Notify(notification)
	}

}
func (d *DefaultGamestate) Disarm(c cards.Card) {
	f := c.GetCost()
	res := d.currentResource
	d.mutex.Lock()
	enough := (&f).IsEnough(res)
	d.mutex.Unlock()
	if enough {
		d.mutex.Lock()
		d.PayResource(f)
		d.mutex.Unlock()
		c.(cards.Trapper).OnDisarm()
		d.RemoveCardFromCenterRow(c)
		// remove c from center cards
		d.updateCenterCard(c)
		// d.BanishCard(c, cards.DISCARD_SOURCE_CENTER)
		c.Dispose(cards.DISCARD_SOURCE_CENTER)
		trapRemovedEvent := map[string]interface{}{cards.EVENT_ATTR_TRAP_REMOVED: c}
		d.NotifyListener(cards.EVENT_TRAP_REMOVED, trapRemovedEvent)

	}

}
func (d *DefaultGamestate) DetachCard(c cards.Overlay) {
	f := c.GetCost()
	res := d.currentResource
	d.MutexLock()
	isEnough := (&f).IsEnough(res)
	d.MutexUnlock()
	if isEnough {
		d.mutex.Lock()
		d.PayResource(f)
		d.mutex.Unlock()
		if c.HasOverlayCard() {
			c.Detach()
		} else {
			c.Dispose(cards.DISCARD_SOURCE_CENTER)
			d.RemoveCardFromCenterRow(c)
			// remove c from center cards
			d.updateCenterCard(c)
		}

	}
}
func (d *DefaultGamestate) DefeatCard(c cards.Card) {
	if !d.LegalCheck(cards.ACTION_DEFEAT, c) {
		return
	}
	f := c.GetCost()
	res := d.currentResource
	d.MutexLock()
	isEnough := (&f).IsEnough(res)
	d.MutexUnlock()
	if isEnough {
		d.mutex.Lock()
		d.PayResource(f)
		d.mutex.Unlock()
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

}
func (d *DefaultGamestate) GetCurrentResource() cards.Resource {
	return d.currentResource
}
func (d *DefaultGamestate) AddResource(name string, amount int) {
	// d.mutex.Lock()
	if amount > 0 {
		d.currentResource.AddResource(name, amount)
	} else {
		d.currentResource.RemoveResource(name, -amount)
	}
	// d.mutex.Unlock()

	if _, ok := d.TopicsListeners[cards.EVENT_ADD_RESOURCE]; ok {
		j := d.TopicsListeners[cards.EVENT_ADD_RESOURCE]
		currVal := d.currentResource.Detail[name]
		data := map[string]interface{}{cards.EVENT_ATTR_ADD_RESOURCE_NAME: name, cards.EVENT_ATTR_ADD_RESOURCE_AMOUNT: currVal}
		j.Notify(data)
	}
}
