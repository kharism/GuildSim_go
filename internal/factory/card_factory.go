package factory

import (
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/cards/item"
	"math/rand"
)

const SET_STARTER_DECK = "starter_deck"

const SET_CENTER_DECK_1 = "center_deck"
const SET_CENTER_DECK_2 = "center_deck_2"

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
		return CreateFillerCenterDeck2(gamestate)
	case SET_CENTER_DECK_2:
		return createStarterCenterDeckAct2(gamestate)
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

// before defeating 1st boss filler
func CreateFillerCenterDeck1(gamestate cards.AbstractGamestate) []cards.Card {
	deck := []cards.Card{}
	for i := 0; i < 3; i++ {
		h := cards.NewGoblinMonster(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 1; i++ {
		h := cards.NewPyroKnight(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 1; i++ {
		h := cards.NewAggroDjinn(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 1; i++ {
		h := cards.NewNoviceAdventurer(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 1; i++ {
		h := cards.NewNoviceCombatant(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 1; i++ {
		h := cards.NewRookieMage(gamestate)
		deck = append(deck, &h)
	}
	return deck
}

// after defeating 1st boss fillter
func CreateFillerCenterDeck2(gamestate cards.AbstractGamestate) []cards.Card {
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
func createStarterCenterDeckAct2(gamestate cards.AbstractGamestate) []cards.Card {
	deck := []cards.Card{}
	for i := 0; i < 4; i++ {
		h := cards.NewTorchtail(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 3; i++ {
		h := cards.NewFirelake(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 6; i++ {
		h := cards.NewAggroDjinn(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 4; i++ {
		h := cards.NewNoviceAdventurer(gamestate)
		deck = append(deck, &h)
		ix := cards.NewNoviceCombatant(gamestate)
		deck = append(deck, &ix)
		j := cards.NewNoviceNurse(gamestate)
		deck = append(deck, &j)
	}
	for i := 0; i < 2; i++ {
		h := cards.NewLightingStag(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 3; i++ {
		h := cards.NewWolfShaman(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 6; i++ {
		h := cards.NewWildBoar(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 6; i++ {
		h := cards.NewDeadweight(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 6; i++ {
		h := cards.NewArcher(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 6; i++ {
		h := cards.NewElephantDjinn(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 6; i++ {
		h := cards.NewWishyDjinn(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 0; i++ {
		h := cards.NewDartTrap(gamestate)
		deck = append(deck, &h)
	}
	for i := 0; i < 1; i++ {
		h := cards.NewSlowTrap(gamestate)
		deck = append(deck, &h)
	}
	// for i := 0; i < 1; i++ {
	// 	h := cards.NewTigerRevenger(gamestate)
	// 	deck = append(deck, &h)
	// }
	// for i := 0; i < 1; i++ {
	// 	h := cards.NewWolfPack(gamestate)
	// 	deck = append(deck, &h)
	// }
	for i := 0; i < 1; i++ {
		h := cards.NewIceWyvern(gamestate)
		deck = append(deck, &h)
	}
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	deck = deck[:20]
	UniqueCard := cards.GetPool(gamestate).Fetch()
	deck = append(deck, UniqueCard)
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}
func createStarterCenterDeck(gamestate cards.AbstractGamestate) []cards.Card {
	deck := []cards.Card{}
	for i := 0; i < 9; i++ {
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
	for i := 0; i < 3; i++ {
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
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	deck = deck[:25]
	UniqueCard := cards.GetPool(gamestate).Fetch()
	deck = append(deck, UniqueCard)
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}
