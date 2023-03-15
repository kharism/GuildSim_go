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
	for i := 0; i < 16; i++ {
		h := cards.NewGoblinMonster(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 5; i++ {
		h := cards.NewGoblinSmallLairArea(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 5; i++ {
		h := cards.NewRookieNurse(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 5; i++ {
		h := cards.NewRookieMage(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 5; i++ {
		h := cards.NewScout(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 5; i++ {
		h := cards.NewMonsterSlayer(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 6; i++ {
		h := cards.NewStagShaman(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 5; i++ {
		h := cards.NewWingedLion(gamestate)
		deck = append(deck, &h)
	}

	return deck
}
