package cards

import (
	"math/rand"
	"time"
)

type CardType int64

const (
	Hero CardType = iota
	Area
	Monster
	Event
	Curse
)

func (i CardType) String() string {
	switch i {
	case Hero:
		return "Hero"
	case Area:
		return "Area"
	case Monster:
		return "Mons"
	case Event:
		return "Event"
	case Curse:
		return "Curse"
	}
	return ""
}

type Card interface {
	GetName() string
	GetDescription() string
	GetCost() Cost

	// return Hero,Area,Monster,Event etc
	GetCardType() CardType

	// when played from hand, do this
	OnPlay()

	// when explored, do this
	OnExplored()

	// when slain, do this
	OnSlain()

	// when discarded to cooldown pile, do this
	OnDiscarded()

	// when added to hand do this
	OnAddedToHand()

	// get rid of this card, you either send this to discard pile or banished pile
	Dispose()
}

func RemoveCard(cards []Card, card Card) []Card {
	newCards := []Card{}

	for _, v := range cards {
		if v == card {
			continue
		}
		newCards = append(newCards, v)
	}
	return newCards
}

type Deck struct {
	cards []Card
}

func (d *Deck) Size() int {
	return len(d.cards)
}

func (d *Deck) Shuffle() {
	shuffler := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffler.Shuffle(d.Size(), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}

// put card on the bottom of deck
func (d *Deck) Push(c Card) {
	d.cards = append(d.cards, c)
}

// put card on top of deck
func (d *Deck) Stack(c Card) {
	l := []Card{c}
	l = append(l, d.cards...)
	d.cards = l
	return
}

func (d *Deck) List() []Card {
	return d.cards
}

func (d *Deck) SetList(l []Card) {
	d.cards = l
}

func (d *Deck) Draw() Card {
	c := d.cards[0]
	j := d.cards[1:]
	d.cards = j
	return c
}
