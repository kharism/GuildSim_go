package cards_test

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/cards/item"
	"github/kharism/GuildSim_go/internal/observer"
	"math/rand"
	"sync"
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
	ChooseMethod     func() int
	ChooseMethodBool func() bool
}

func (t *TestCardPicker) ShowDetail(cards.Card) {

}
func (t *TestCardPicker) PickCardOptional(list []cards.Card, message string) int {
	return t.ChooseMethod()
}
func (t *TestCardPicker) BoolPick(message string) bool {
	return t.ChooseMethodBool()
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
	CardsInDeck       cards.DeterministicDeck
	CardsInCenterDeck cards.DeterministicDeck
	TopicsListeners   map[string]*DummyEventListener
	RuleEnforcer      map[string]*cards.RuleEnforcer
	CardsInHand       []cards.Card
	CardsPlayed       []cards.Card
	CenterCards       []cards.Card
	CardsBanished     []cards.Card
	ItemCards         []cards.Card
	CardsDiscarded    cards.DeterministicDeck
	HitPoint          int
	mutex             sync.Mutex
	//ui stuff
	cardPiker  cards.AbstractCardPicker
	boolPiker  cards.AbstractBoolPicker
	cardViewer cards.AbstractDetailViewer
}

func (d *DummyGamestate) MutexLock() {
	d.mutex.Lock()
}
func (d *DummyGamestate) MutexUnlock() {
	d.mutex.Unlock()
}
func (d *DummyGamestate) PayResource(cost cards.Cost) {
	for key, val := range cost.Detail {
		d.currentResource.Detail[key] -= val
	}
}
func (d *DummyGamestate) SetFillerIndex(i int) {}
func (d *DummyGamestate) GetFillerIndex() int  { return 0 }
func (d *DummyGamestate) AppendCardFiller(func(cards.AbstractGamestate) []cards.Card) {

}
func (d *DummyGamestate) BeginTurn() {
	//d.centerCardChanged = false
	for i := 0; i < 5; i++ {
		d.Draw()
	}
	if _, ok := d.TopicsListeners[cards.EVENT_START_OF_TURN]; ok {
		j := d.TopicsListeners[cards.EVENT_START_OF_TURN]
		data := map[string]interface{}{}
		j.Notify(data)
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
	d.CardsInDeck = cards.DeterministicDeck{}
	d.CardsDiscarded = cards.DeterministicDeck{}
	d.CardsBanished = []cards.Card{}
	d.ItemCards = []cards.Card{}
	d.RuleEnforcer = map[string]*cards.RuleEnforcer{}
	d.mutex = sync.Mutex{}
	d.HitPoint = 60
	return &d
}
func (d *DummyGamestate) GetCurrentHP() int {
	return d.HitPoint
}
func (d *DummyGamestate) TakeDamage(dmg int) {
	if dmg > 0 {
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
func (d *DummyGamestate) SetDetailViewer(v cards.AbstractDetailViewer) {
	d.cardViewer = v
}
func (d *DummyGamestate) GetDetailViewer() cards.AbstractDetailViewer {
	return d.cardViewer
}
func (d *DummyGamestate) PeekCenterCard() cards.Card {
	if len(d.CardsInCenterDeck.List()) == 0 {
		return nil
	}
	return d.CardsInCenterDeck.List()[0]
}
func (d *DummyGamestate) GenerateRandomPotion(rarity int) cards.Card {
	return item.CreatePotionRandom(d, rarity)
}
func (d *DummyGamestate) GenerateRandomRelic(rarity int) cards.Card {
	return item.CreateRelicRandom(d, rarity)
}
func (d *DummyGamestate) NotifyListener(eventname string, data map[string]interface{}) {
	if _, ok := d.TopicsListeners[eventname]; ok {
		j := d.TopicsListeners[eventname]
		j.Notify(data)
	}
}
func (d *DummyGamestate) GetMainDeck() *cards.Deck {
	return &d.CardsInDeck.Deck
}
func (d *DummyGamestate) ActDecorators() []func(cards.AbstractGamestate) cards.AbstractGamestate {
	return nil
}
func (d *DummyGamestate) AddActDecorator(f func(cards.AbstractGamestate) cards.AbstractGamestate) {
	return
}
func (d *DummyGamestate) GetCardPicker() cards.AbstractCardPicker {
	return d.cardPiker
}
func (d *DummyGamestate) SetCardPicker(a cards.AbstractCardPicker) {
	d.cardPiker = a
}
func (d *DummyGamestate) GetBoolPicker() cards.AbstractBoolPicker {
	return d.boolPiker
}
func (d *DummyGamestate) SetBoolPicker(a cards.AbstractBoolPicker) {
	d.boolPiker = a
}
func (d *DummyGamestate) AddQuest(s string)    {}
func (d *DummyGamestate) RemoveQuest(s string) {}
func (d *DummyGamestate) GetQuests() []string {
	return []string{}
}
func (d *DummyGamestate) EndTurn() {
	// reset resource except money and reputation
	curRes := d.GetCurrentResource().Detail
	for k := range curRes {
		if k == cards.RESOURCE_NAME_MONEY || k == cards.RESOURCE_NAME_REPUTATION || k == cards.RESOURCE_NAME_BLOCK {
			continue
		}
		d.GetCurrentResource().Detail[k] = 0
	}

	// remove cards played
	for _, c := range d.CardsPlayed {
		// d.CardsDiscarded.Push(c)
		c.Dispose(cards.DISCARD_SOURCE_PLAYED)
		if pun, ok := c.(cards.Punisher); ok {
			pun.OnPunish()
		}
	}
	d.CardsPlayed = []cards.Card{}

	// remove cards in hand
	for _, c := range d.CardsInHand {
		// d.CardsDiscarded.Push(c)
		c.Dispose(cards.DISCARD_SOURCE_HAND)
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
	d.GetCurrentResource().Detail[cards.RESOURCE_NAME_BLOCK] = 0
	if _, ok := d.TopicsListeners[cards.EVENT_END_OF_TURN]; ok {
		j := d.TopicsListeners[cards.EVENT_END_OF_TURN]
		data := map[string]interface{}{}
		j.Notify(data)
	}
}
func (d *DummyGamestate) AddCardToCenterDeck(source string, shuffle bool, c ...cards.Card) {
	for _, cc := range c {
		d.CardsInCenterDeck.Stack(cc)
	}
	if shuffle {
		d.CardsInCenterDeck.Shuffle()
	}

	// fmt.Println("Done adding", d.CardsInCenterDeck.Size(), "To center deck")
}

func (d *DummyGamestate) StackCards(source string, cc ...cards.Card) {
	for _, c := range cc {
		d.CardsInDeck.Stack(c)
		c.OnReturnToDeck()
		if _, ok := d.TopicsListeners[cards.EVENT_ATTR_CARD_STACKED]; ok {
			j := d.TopicsListeners[cards.EVENT_ATTR_CARD_STACKED]
			data := map[string]interface{}{cards.EVENT_ATTR_CARD_STACKED: c}
			j.Notify(data)
		}
	}

}
func (d *DummyGamestate) ShuffleMainDeck() {
	d.CardsInDeck.Shuffle()
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
	if !d.LegalCheck(cards.ACTION_RECRUIT, c) {
		return
	}
	k := c.GetCost()
	if k.IsEnough(d.currentResource) {
		d.PayResource(k)
		if _, ok := c.(cards.Recruitable); ok {
			f := c.(cards.Recruitable)
			f.OnRecruit()
		}
		// d.DiscardCard(c)
		d.CardsDiscarded.Push(c)
		cardRecruitedEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_RECRUITED: c}
		l, ok := d.TopicsListeners[cards.EVENT_CARD_RECRUITED]
		if ok {
			l.Notify(cardRecruitedEvent)
		}
		// remove c from center cards
		d.updateCenterCard(c)
	}
	return
}
func (d *DummyGamestate) GetCooldownCard() []cards.Card {
	return d.CardsDiscarded.List()
}
func (d *DummyGamestate) DiscardCard(c cards.Card, source string) {
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
func (d *DummyGamestate) RemoveCardFromPlayed(c cards.Card) {
	for idx, c2 := range d.CardsPlayed {
		if c2 == c {
			d.RemoveCardFromHandIdx(idx)
			return
		}
	}
}
func (d *DummyGamestate) RemoveCardFromPlayedIdx(i int) {
	j := append(d.CardsPlayed[:i], d.CardsPlayed[i+1:]...)
	d.CardsPlayed = j
}
func (d *DummyGamestate) DetachCard(c cards.Overlay) {
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
func (d *DummyGamestate) AttachLegalCheck(actionName string, lc cards.LegalChecker) {
	if _, ok := d.RuleEnforcer[actionName]; !ok {
		d.RuleEnforcer[actionName] = cards.NewRuleEnforcer()
	}
	d.RuleEnforcer[actionName].AttachRule(lc)
}
func (d *DummyGamestate) DetachLegalCheck(actionName string, lc cards.LegalChecker) {
	if _, ok := d.RuleEnforcer[actionName]; !ok {
		return
	}
	d.RuleEnforcer[actionName].DetachRule(lc)
}
func (d *DummyGamestate) LegalCheck(actionName string, data interface{}) bool {
	j, ok := d.RuleEnforcer[actionName]
	if !ok {
		return true
	}
	return j.Check(data)
}
func (d *DummyGamestate) RemoveCardFromCooldown(c cards.Card) {
	for idx, c2 := range d.CardsDiscarded.List() {
		if c2 == c {
			d.RemoveCardFromCooldownIdx(idx)
			return
		}
	}
}
func (d *DummyGamestate) AppendCenterCard(c cards.Card) {
	d.CenterCards = append(d.CenterCards, c)
}
func (d *DummyGamestate) RemoveCardFromCooldownIdx(i int) {
	cooldownList := d.CardsDiscarded.List()
	j := append(cooldownList[:i], cooldownList[i+1:]...)
	d.CardsDiscarded.SetList(j)
}
func (d *DummyGamestate) UpdateCenterCard(c cards.Card) {
	d.updateCenterCard(c)
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
	if !d.LegalCheck(cards.ACTION_EXPLORE, c) {
		return
	}
	// check cost and resource
	f := c.GetCost()
	res := d.currentResource
	if (&f).IsEnough(res) {
		// payResource
		d.PayResource(f)
		c.OnExplored()
		d.BanishCard(c, cards.DISCARD_SOURCE_CENTER)
		cardExploredEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_EXPLORED: c}

		l, ok := d.TopicsListeners[cards.EVENT_CARD_EXPLORED]
		if ok {
			l.Notify(cardExploredEvent)
		}
		// remove c from center cards
		d.updateCenterCard(c)
	}
}
func (d *DummyGamestate) Disarm(c cards.Card) {
	if !d.LegalCheck(cards.ACTION_DISARM, c) {
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
		trapRemovedEvent := map[string]interface{}{cards.EVENT_ATTR_TRAP_REMOVED: c}
		d.NotifyListener(cards.EVENT_TRAP_REMOVED, trapRemovedEvent)

	}
	return
}
func (d *DummyGamestate) ReplaceCenterCard() cards.Card {
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
		}
	}
	return replacementCard
}
func (d *DummyGamestate) Draw() {
	if !d.LegalCheck(cards.ACTION_DRAW, nil) {
		return
	}
	if d.CardsInDeck.Size() == 0 {
		// shuffle discard pile
		d.CardsDiscarded.Shuffle()
		d.CardsInDeck = d.CardsDiscarded
		d.CardsDiscarded = cards.DeterministicDeck{}
	}
	newCard := d.CardsInDeck.Draw()
	d.CardsInHand = append(d.CardsInHand, newCard)
	newCard.OnAddedToHand()
	return
}
func (d *DummyGamestate) BanishCard(c cards.Card, source string) {
	d.CardsBanished = append(d.CardsBanished, c)
	c.OnBanished()
	if _, ok := d.TopicsListeners[cards.EVENT_CARD_BANISHED]; ok {
		notification := map[string]interface{}{cards.EVENT_ATTR_CARD_BANISHED: c, cards.EVENT_ATTR_DISCARD_SOURCE: source}
		d.TopicsListeners[cards.EVENT_CARD_BANISHED].Notify(notification)
	}
	return
}
func (d *DummyGamestate) DefeatCard(c cards.Card) {
	if !d.LegalCheck(cards.ACTION_DEFEAT, c) {
		return
	}
	fmt.Println("Trying to defeat")
	f := c.GetCost()
	res := d.currentResource
	if (&f).IsEnough(res) {
		d.PayResource(f)
		c.OnSlain()
		// d.BanishCard(c, cards.DISCARD_SOURCE_CENTER)
		d.CardsBanished = append(d.CardsBanished, c)
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
func (d *DummyGamestate) AddItem(c cards.Card) {
	d.ItemCards = append(d.ItemCards, c)
	c.OnAcquire()
}
func (d *DummyGamestate) ConsumeItem(c cards.Consumable) {
	c.OnConsume()
}
func (d *DummyGamestate) ListItems() []cards.Card {
	return d.ItemCards
}
func (d *DummyGamestate) RemoveItem(c cards.Card) {
	idx := -1
	for i := range d.ItemCards {
		if d.ItemCards[i] == c {
			idx = i
			break
		}
	}
	if idx != -1 {
		d.RemoveItemIndex(idx)
	}
}
func (d *DummyGamestate) RemoveItemIndex(i int) {
	h := append(d.ItemCards[:i], d.ItemCards[i+1:]...)
	d.ItemCards = h
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
