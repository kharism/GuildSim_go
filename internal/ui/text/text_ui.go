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

type TextCardPicker struct {
	// ChooseMethod func() int
}

func (t *TextCardPicker) PickCardOptional(list []cards.Card, message string) int {
	fmt.Println(message)
	for i, card := range list {
		fmt.Printf("[%d] %s [%s]\n", i, card.GetName(), card.GetCost())
	}
	// reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(os.Stdin)

	// text, _ := reader.ReadString('\n')
	for scanner.Scan() {
		if scanner.Text() == "" {
			return -1
		}
		picks, err := strconv.Atoi(scanner.Text())
		if err != nil || picks < 0 || picks > len(list) {
			continue
		}
		return picks
	}
	return -1
	// return t.ChooseMethod()
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
		if err != nil || picks < 0 || picks > len(list) {
			continue
		}
		return picks
	}
	return -1
	// return t.ChooseMethod()
}

type TextUIGamestate struct {
	gamestate  cards.AbstractGamestate
	cardPicker cards.AbstractCardPicker
}

func NewTextUIGamestate(gamestate cards.AbstractGamestate) *TextUIGamestate {
	d := TextUIGamestate{}
	d.gamestate = gamestate
	d.cardPicker = &TextCardPicker{}
	// d.HitPoint = 60
	return &d
}
func (d *TextUIGamestate) GetCooldownCard() []cards.Card {
	return d.gamestate.GetCooldownCard()
}
func (d *TextUIGamestate) AddCardToCenterDeck(source string, c ...cards.Card) {
	d.AddCardToCenterDeck(source, c...)
}

func (d *TextUIGamestate) PayResource(cost cards.Cost) {
	// for key, val := range cost.Detail {
	// 	d.currentResource.Detail[key] -= val
	// }
	d.gamestate.PayResource(cost)
}

func (d *TextUIGamestate) AttachListener(eventName string, l observer.Listener) {
	// if _, ok := d.TopicsListeners[eventName]; !ok {
	// 	d.TopicsListeners[eventName] = &DummyEventListener{}
	// }
	// k := (d.TopicsListeners[eventName])
	// k.Attach(l)
	d.gamestate.AttachListener(eventName, l)
	// fmt.Println("Attach Listener", len(d.TopicsListeners[eventName].Listeners))
}
func (d *TextUIGamestate) RemoveListener(eventName string, l observer.Listener) {
	// if _, ok := d.TopicsListeners[eventName]; !ok {
	// 	return
	// }
	// // fmt.Println("Remove Listener")
	// k := (d.TopicsListeners[eventName])
	// k.Detach(l)
	d.gamestate.RemoveListener(eventName, l)
}

func (d *TextUIGamestate) GetCurrentHP() int {
	// return d.HitPoint
	return d.GetCurrentHP()
}
func (d *TextUIGamestate) TakeDamage(dmg int) {
	// d.HitPoint -= dmg
	// l, ok := d.TopicsListeners[cards.EVENT_TAKE_DAMAGE]
	// takeDamageEvent := map[string]interface{}{cards.EVENT_TAKE_DAMAGE: dmg}
	// if ok {
	// 	l.Notify(takeDamageEvent)
	// }
	d.gamestate.TakeDamage(dmg)
}
func (d *TextUIGamestate) GetCardPicker() cards.AbstractCardPicker {
	return d.cardPicker
}
func (d *TextUIGamestate) EndTurn() {
	d.gamestate.EndTurn()
	for i := 0; i < 5; i++ {
		d.Draw()
	}
}
func (d *TextUIGamestate) PlayCard(c cards.Card) {
	d.gamestate.PlayCard(c)
}
func (d *TextUIGamestate) GetPlayedCards() []cards.Card {
	return d.gamestate.GetPlayedCards()
}
func (d *TextUIGamestate) GetCardInHand() []cards.Card {
	return d.gamestate.GetCardInHand()
}
func (d *TextUIGamestate) GetCenterCard() []cards.Card {
	return d.gamestate.GetCenterCard()
}
func (d *TextUIGamestate) RecruitCard(c cards.Card) {
	d.gamestate.RecruitCard(c)
	return
}

// func (d *TextUIGamestate) DiscardCard(c cards.Card) {
// 	d.gamestate.DiscardCard(c)
// 	return
// }
func (d *TextUIGamestate) CenterRowInit() {
	d.gamestate.CenterRowInit()
}

func (d *TextUIGamestate) Explore(c cards.Card) {
	d.gamestate.Explore(c)
}
func (d *TextUIGamestate) ReplaceCenterCard() cards.Card {
	return d.gamestate.ReplaceCenterCard()
}
func (d *TextUIGamestate) Draw() {
	d.gamestate.Draw()
}

// func (d *TextUIGamestate) BanishCard(c cards.Card) {
// 	d.gamestate.BanishCard(c)
// 	return
// }
func (d *TextUIGamestate) DefeatCard(c cards.Card) {
	d.gamestate.DefeatCard(c)
}
func (d *TextUIGamestate) GetCurrentResource() cards.Resource {
	return d.gamestate.GetCurrentResource()
}
func (d *TextUIGamestate) AddResource(name string, amount int) {
	d.gamestate.AddResource(name, amount)
}
func (d *TextUIGamestate) RemoveCardFromHand(c cards.Card) {
	d.gamestate.RemoveCardFromHand(c)
}
func (d *TextUIGamestate) RemoveCardFromHandIdx(i int) {
	d.gamestate.RemoveCardFromHandIdx(i)
}
func (d *TextUIGamestate) RemoveCardFromCenterRow(c cards.Card) {
	d.gamestate.RemoveCardFromCenterRow(c)
}
func (d *TextUIGamestate) RemoveCardFromCenterRowIdx(i int) {
	d.gamestate.RemoveCardFromCenterRowIdx(i)
}
func (d *TextUIGamestate) RemoveCardFromCooldown(c cards.Card) {
	d.gamestate.RemoveCardFromCooldown(c)
}
func (d *TextUIGamestate) RemoveCardFromCooldownIdx(idx int) {
	d.gamestate.RemoveCardFromCooldownIdx(idx)
}
func InverseCenterCardsKey(s string) int {
	i := "qwert"
	for idx, v := range i {
		if strings.ToLower(string(v)) == strings.ToLower(s) {
			fmt.Println("access center card idx", idx)
			return idx
		}
	}
	return -1
}
func SwitchCenterCardsKey(idx int, card cards.Card, resource *cards.Resource) string {
	cost := card.GetCost()
	if !cost.IsEnough(*resource) {
		return ""
	}
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
	orderOfResource := []string{cards.RESOURCE_NAME_COMBAT, cards.RESOURCE_NAME_EXPLORATION, cards.RESOURCE_NAME_REPUTATION}
	ResourceStr := []string{fmt.Sprintf("HP: %d", d.gamestate.GetCurrentHP())}
	// for key, val := range d.gamestate.GetCurrentResource().Detail {
	// 	// fmt.Printf("%s:%d  \n", key, val)
	// 	l := fmt.Sprintf("%s:%d", key, val)
	// 	ResourceStr = append(ResourceStr, l)
	// }
	res := d.gamestate.GetCurrentResource()
	for _, val := range orderOfResource {
		if _, ok := res.Detail[val]; ok {
			l := fmt.Sprintf("%s:%d", val, res.Detail[val])
			ResourceStr = append(ResourceStr, l)
		}

	}
	fmt.Println(strings.Join(ResourceStr, " "))
	fmt.Println("Cards In center Row:")
	for idx, card := range d.gamestate.GetCenterCard() {
		fmt.Printf("[%s] %s (%s) [%s]:%s\n", SwitchCenterCardsKey(idx, card, &res), card.GetName(), card.GetCardType(), card.GetCost(), card.GetDescription())
	}
	fmt.Println("Cards in hand:")
	for idx, card := range d.gamestate.GetCardInHand() {
		fmt.Printf("[%d] %s [%s]:%s\n", idx, card.GetName(), card.GetCost(), card.GetDescription())
	}
}
func (d *TextUIGamestate) Run() {
	// initialize gamestate
	d.gamestate.CenterRowInit()
	//draw 6

	for i := 0; i < 5; i++ {
		d.gamestate.Draw()
	}
	scanner := bufio.NewScanner(os.Stdin)
	d.Render()
	// text, _ := reader.ReadString('\n')
	for scanner.Scan() {
		input := scanner.Text()
		ll, err := strconv.Atoi(input)
		if err != nil {
			// NaN meaning we acquire cards or end our turn
			if input == "" {
				d.EndTurn()
				d.Render()
				continue
			} else if input == "A" {
				cardsInHand := d.gamestate.GetCardInHand()
				for i := len(cardsInHand) - 1; i >= 0; i-- {
					d.gamestate.PlayCard(cardsInHand[i])
				}
				d.Render()
				continue
			}
			cardIdx := InverseCenterCardsKey(input)
			choosenCard := d.gamestate.GetCenterCard()[cardIdx]
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
			if ll < 0 || ll >= len(d.gamestate.GetCardInHand()) {
				fmt.Println("invalid number, try again")
				continue
			}
			choosenCard := d.gamestate.GetCardInHand()[ll]
			d.PlayCard(choosenCard)
		}
		d.Render()
	}
}
