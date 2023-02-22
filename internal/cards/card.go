package cards

import "math/rand"

type Card interface {
	GetName() string
	GetDescription() string
	GetCost() Cost

	// when played on hand, to this
	OnPlay()
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
