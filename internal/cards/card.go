package cards

import "math/rand"

type CardType int64

const (
	Hero CardType = iota
	Area
	Monster
	Event
)

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
	rand.Shuffle(d.Size(), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
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

func (d *Deck) Draw() Card {
	c := d.cards[0]
	j := d.cards[1:]
	d.cards = j
	return c
}
