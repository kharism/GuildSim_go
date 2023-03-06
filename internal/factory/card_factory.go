package factory

import "github/kharism/GuildSim_go/internal/cards"

func CreateStarterDeck(gamestate cards.AbstractGamestate) []cards.Card {
	deck := []cards.Card{}
	for i := 0; i < 5; i++ {
		h := cards.NewRookieAdventurer(gamestate)
		deck = append(deck, &h)
		j := cards.NewRookieCombatant(gamestate)
		deck = append(deck, &j)
	}
	return deck
}

func CreateStarterCenterDeck(gamestate cards.AbstractGamestate) []cards.Card {
	deck := []cards.Card{}
	for i := 0; i < 4; i++ {
		h := cards.NewGoblinMonster(gamestate)
		deck = append(deck, &h)
	}
	return deck
}
