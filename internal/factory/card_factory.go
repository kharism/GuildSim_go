package factory

import (
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/cards/item"
)

const SET_STARTER_DECK = "starter_deck"

const SET_CENTER_DECK_1 = "center_deck"

const SET_POTION_COMMON_RANDOM = "POT_COMMON_RAND"

const SET_FILLER_CARDS = "filler"

// generate list of cards. use public constant to pick which set of cards to generate
// for future devs: add in more functions to create expansion set that just need to be slapped in without adding triggers/eventListeners
// if you want to add expansion with event listeners use decorators
func CardFactory(setname string, gamestate cards.AbstractGamestate) []cards.Card {
	switch setname {
	case SET_STARTER_DECK:
		return createStarterDeck(gamestate)
	case SET_POTION_COMMON_RANDOM:
		item.CreatePotionRandom(gamestate, cards.RARITY_COMMON)
	case SET_FILLER_CARDS:
		return createFillerCenterDeck(gamestate)
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
func createFillerCenterDeck(gamestate cards.AbstractGamestate) []cards.Card {
	deck := []cards.Card{}
	for i := 0; i < 3; i++ {
		h := cards.NewGoblinSmallLairArea(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 4; i++ {
		h := cards.NewIceWyvern(gamestate)
		j := cards.NewTorchtail(gamestate)
		deck = append(deck, &h, &j)
	}
	return deck
}
func createStarterCenterDeck(gamestate cards.AbstractGamestate) []cards.Card {
	deck := []cards.Card{}
	for i := 0; i < 12; i++ {
		h := cards.NewGoblinMonster(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 6; i++ {
		h := cards.NewWildBoar(gamestate)
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
	for i := 0; i < 3; i++ {
		h := cards.NewDeadweight(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 2; i++ {
		h := cards.NewMonsterSlayer(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 2; i++ {
		h := cards.NewCleric(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 3; i++ {
		ll := cards.NewBulwark(gamestate)
		deck = append(deck, &ll)
	}
	ll := cards.NewArcher(gamestate)
	deck = append(deck, &ll)
	for i := 0; i < 2; i++ {
		h := cards.NewRookieHunter(gamestate)
		j := cards.NewThief(gamestate)
		deck = append(deck, &h, &j)
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
