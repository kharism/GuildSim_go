package factory

import "github/kharism/GuildSim_go/internal/cards"

const SET_STARTER_DECK = "starter_deck"

const SET_CENTER_DECK_1 = "center_deck"

// generate list of cards. use public constant to pick which set of cards to generate
// for future devs: add in more functions to create expansion set that just need to be slapped in without adding triggers/eventListeners
// if you want to add expansion with event listeners use decorators
func CardFactory(setname string, gamestate cards.AbstractGamestate) []cards.Card {
	switch setname {
	case SET_STARTER_DECK:
		return createStarterDeck(gamestate)
	case SET_CENTER_DECK_1:
	default:
		return createStarterCenterDeck(gamestate)

	}
	return createStarterCenterDeck(gamestate)
}

func createStarterDeck(gamestate cards.AbstractGamestate) []cards.Card {
	deck := []cards.Card{}
	for i := 0; i < 5; i++ {
		h := cards.NewRookieAdventurer(gamestate)
		deck = append(deck, &h)
		j := cards.NewRookieCombatant(gamestate)
		deck = append(deck, &j)
	}
	return deck
}

func createStarterCenterDeck(gamestate cards.AbstractGamestate) []cards.Card {
	deck := []cards.Card{}
	for i := 0; i < 12; i++ {
		h := cards.NewGoblinMonster(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 3; i++ {
		h := cards.NewGoblinSmallLairArea(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 3; i++ {
		h := cards.NewRookieNurse(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 3; i++ {
		h := cards.NewRookieMage(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 3; i++ {
		h := cards.NewScout(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 3; i++ {
		h := cards.NewPyroKnight(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 3; i++ {
		h := cards.NewFireMage(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 2; i++ {
		h := cards.NewMonsterSlayer(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 2; i++ {
		h := cards.NewStagShaman(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 2; i++ {
		h := cards.NewWingedLion(gamestate)
		deck = append(deck, &h)
	}

	return deck
}
