package gamestate

import (
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/decorator"
	"github/kharism/GuildSim_go/internal/factory"
)

// create customized DefaultGamestate using several parameters
// starterDeckSets and centerDecksets refers to string const in factory.
// decorators is from decorator package
func CustomizedDefaultGamestate(starterDeckSets, centerDeckSets []string, decorators []decorator.AbstractDecorator) cards.AbstractGamestate {
	defGameState := DefaultGamestate{}
	defGameState.cardPiker = nil
	defGameState.currentResource = cards.NewResource()
	defGameState.CardsPlayed = []cards.Card{}
	defGameState.TopicsListeners = map[string]*DummyEventListener{}
	defGameState.CenterCards = []cards.Card{}
	defGameState.CardsInHand = []cards.Card{}
	// d.cardPiker = &TextCardPicker{}
	defGameState.HitPoint = 60
	defGameState.CardsDiscarded = cards.Deck{}
	defGameState.CardsInCenterDeck = cards.Deck{}
	centerDeck := []cards.Card{}
	for _, setName := range centerDeckSets {
		newCards := factory.CardFactory(setName, &defGameState)
		centerDeck = append(centerDeck, newCards...)
	}
	defGameState.CardsInCenterDeck.SetList(centerDeck)
	defGameState.CardsInDeck = cards.Deck{}
	starterDeck := []cards.Card{}
	for _, setName := range starterDeckSets {
		newCards := factory.CardFactory(setName, &defGameState)
		starterDeck = append(starterDeck, newCards...)
	}
	defGameState.CardsInDeck.SetList(starterDeck)

	var j cards.AbstractGamestate
	// j := decorators[0](&defGameState)
	for _, k := range decorators {
		j = k(&defGameState)
	}
	if j == nil {
		return &defGameState
	} else {
		return j
	}
}
