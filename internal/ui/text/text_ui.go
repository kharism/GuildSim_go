package text

import (
	"bufio"
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/observer"
	"os"
	"strconv"
	"strings"
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

type TextCardPicker struct {
	// ChooseMethod func() int
}

func (t *TextCardPicker) PickCard(list []cards.Card, message string) int {
	fmt.Println(message)
	for i, card := range list {
		fmt.Printf("[%d] %s [%s]\n", i, card.GetName(), card.GetCost())
	}
	// reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(os.Stdin)

	// text, _ := reader.ReadString('\n')
	for scanner.Scan() {
		picks, err := strconv.Atoi(scanner.Text())
		if err != nil {
			continue
		}
		return picks
	}
	return -1
	// return t.ChooseMethod()
}

type TextUIGamestate struct {
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

func NewTextUIGamestate() cards.AbstractGamestate {
	d := TextUIGamestate{}
	d.currentResource = cards.NewResource()
	d.CardsPlayed = []cards.Card{}
	d.TopicsListeners = map[string]*DummyEventListener{}
	d.CenterCards = []cards.Card{}
	d.CardsInHand = []cards.Card{}
	d.cardPiker = &TextCardPicker{}
	d.HitPoint = 60
	return &d
}

func (d *TextUIGamestate) AddCardToCenterDeck(c ...cards.Card) {
	for _, cc := range c {
		d.CardsInCenterDeck.Stack(cc)
	}
	d.CardsInCenterDeck.Shuffle()
}

func (d *TextUIGamestate) PayResource(cost cards.Cost) {
	for key, val := range cost.Detail {
		d.currentResource.Detail[key] -= val
	}
}

func (d *TextUIGamestate) AttachListener(eventName string, l observer.Listener) {
	if _, ok := d.TopicsListeners[eventName]; !ok {
		d.TopicsListeners[eventName] = &DummyEventListener{}
	}
	k := (d.TopicsListeners[eventName])
	k.Attach(l)
	// fmt.Println("Attach Listener", len(d.TopicsListeners[eventName].Listeners))
}
func (d *TextUIGamestate) RemoveListener(eventName string, l observer.Listener) {
	if _, ok := d.TopicsListeners[eventName]; !ok {
		return
	}
	// fmt.Println("Remove Listener")
	k := (d.TopicsListeners[eventName])
	k.Detach(l)
}

func (d *TextUIGamestate) GetCurrentHP() int {
	return d.HitPoint
}
func (d *TextUIGamestate) TakeDamage(dmg int) {
	d.HitPoint -= dmg
	l, ok := d.TopicsListeners[cards.EVENT_TAKE_DAMAGE]
	takeDamageEvent := map[string]interface{}{cards.EVENT_TAKE_DAMAGE: dmg}
	if ok {
		l.Notify(takeDamageEvent)
	}
}
func (d *TextUIGamestate) GetCardPicker() cards.AbstractCardPicker {
	return d.cardPiker
}
func (d *TextUIGamestate) EndTurn() {
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
func (d *TextUIGamestate) PlayCard(c cards.Card) {
	c.OnPlay()
	// fmt.Println("Card played", c.GetName())
	cardPlayedEvent := map[string]interface{}{cards.EVENT_ATTR_CARD_PLAYED: c}

	l, ok := d.TopicsListeners[cards.EVENT_CARD_PLAYED]
	if ok {
		l.Notify(cardPlayedEvent)
	}

	d.CardsPlayed = append(d.CardsPlayed, c)
}
func (d *TextUIGamestate) GetPlayedCards() []cards.Card {
	return d.CardsPlayed
}
func (d *TextUIGamestate) GetCardInHand() []cards.Card {
	return d.CardsInHand
}
func (d *TextUIGamestate) GetCenterCard() []cards.Card {
	return d.CenterCards
}
func (d *TextUIGamestate) RecruitCard(c cards.Card) {
	return
}
func (d *TextUIGamestate) DiscardCard(c cards.Card) {
	return
}
func (d *TextUIGamestate) CenterRowInit() {
	f := d.ReplaceCenterCard()
	d.CenterCards = append(d.CenterCards, f)
}
func (d *TextUIGamestate) updateCenterCard(c cards.Card) {
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
func (d *TextUIGamestate) Explore(c cards.Card) {
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
func (d *TextUIGamestate) ReplaceCenterCard() cards.Card {
	return d.CardsInCenterDeck.Draw()
}
func (d *TextUIGamestate) Draw() {
	newCard := d.CardsInDeck.Draw()
	d.CardsInHand = append(d.CardsInHand, newCard)
	return
}
func (d *TextUIGamestate) BanishCard(c cards.Card) {
	return
}
func (d *TextUIGamestate) DefeatCard(c cards.Card) {
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
func (d *TextUIGamestate) GetCurrentResource() cards.Resource {
	return d.currentResource
}
func (d *TextUIGamestate) AddResource(name string, amount int) {
	d.currentResource.AddResource(name, amount)
}
func InverseCenterCardsKey(s string) int {
	i := "qwert"
	for idx, v := range i {
		if strings.ToLower(string(v)) == strings.ToLower(s) {
			return idx
		}
	}
	return -1
}
func SwitchCenterCardsKey(idx int) string {
	switch idx {
	case 0:
		return "q"
	case 1:
		return "w"
	case 2:
		return "e"
	case 3:
		return "r"
	case 4:
		return "t"
	}
	return ""
}
func (d *TextUIGamestate) Render() {
	fmt.Println("Resource")
	for key, val := range d.currentResource.Detail {
		fmt.Print("%s:%d  ", key, val)
	}
	fmt.Println("Cards In center Row:")
	for idx, card := range d.CenterCards {
		fmt.Printf("[%s] %s [%s]:%s\n", SwitchCenterCardsKey(idx), card.GetName(), card.GetCost(), card.GetDescription())
	}
	fmt.Println("Cards in hand:")
	for idx, card := range d.CenterCards {
		fmt.Printf("[%d] %s [%s]:%s\n", idx, card.GetName(), card.GetCost(), card.GetDescription())
	}
}
func (d *TextUIGamestate) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	d.Render()
	// text, _ := reader.ReadString('\n')
	for scanner.Scan() {
		input := scanner.Text()
		ll, err := strconv.Atoi(input)
		if err != nil {
			// NaN meaning we acquire cards
			cardIdx := InverseCenterCardsKey(input)
			choosenCard := d.CenterCards[cardIdx]
			switch choosenCard.GetCardType() {
			case cards.Monster:
				d.DefeatCard(choosenCard)
			case cards.Area:
				d.Explore(choosenCard)
			case cards.Hero:
				d.RecruitCard(choosenCard)
			}
		} else {
			// we play a card
			choosenCard := d.CenterCards[ll]
			d.PlayCard(choosenCard)
		}
		d.Render()
	}
}
