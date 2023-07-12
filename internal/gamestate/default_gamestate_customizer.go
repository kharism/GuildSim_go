package gamestate

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/decorator"
	"github/kharism/GuildSim_go/internal/factory"
)

type BossDefeatedAction struct {
	State *DefaultGamestate
}

func (a *BossDefeatedAction) DoAction(data map[string]interface{}) {
	count := data[cards.EVENT_ATTR_BOSS_DEFEATED_COUNT].(int)
	// shuffle all discard and hand to deck
	// remove all curse card and resource in the state
	NewDeck := []cards.Card{}
	cardsInDeck := a.State.CardsInDeck.List()
	cardsDiscarded := a.State.CardsDiscarded.List()
	cardsInHand := a.State.CardsInHand
	for _, c := range cardsInDeck {
		if c.GetCardType() != cards.Curse {
			NewDeck = append(NewDeck, c)
		}
	}
	for _, c := range cardsDiscarded {
		if c.GetCardType() != cards.Curse {
			NewDeck = append(NewDeck, c)
		}
	}
	for _, c := range cardsInHand {
		if c.GetCardType() != cards.Curse {
			NewDeck = append(NewDeck, c)
		}
	}
	for _, c := range a.State.CardsPlayed {
		if c.GetCardType() != cards.Curse {
			NewDeck = append(NewDeck, c)
		}
	}
	a.State.CardsInHand = []cards.Card{}
	blockCount := a.State.currentResource.Detail[cards.RESOURCE_NAME_BLOCK]
	combatCount := a.State.currentResource.Detail[cards.RESOURCE_NAME_COMBAT]
	explorationCount := a.State.currentResource.Detail[cards.RESOURCE_NAME_EXPLORATION]
	a.State.AddResource(cards.RESOURCE_NAME_BLOCK, -blockCount)
	a.State.AddResource(cards.RESOURCE_NAME_COMBAT, -combatCount)
	a.State.AddResource(cards.RESOURCE_NAME_EXPLORATION, -explorationCount)
	a.State.CardsBanished = []cards.Card{}
	a.State.CardsInDeck.SetList(NewDeck)
	// generate new center deck based on which act we've just cleared
	if count == 1 {
		// attach base center deck
		newMainDeckCard := factory.CardFactory(factory.SET_CENTER_DECK_2, a.State)
		a.State.CardsInCenterDeck.SetList(newMainDeckCard)
		a.State.CenterCards = []cards.Card{}
		// attach the main act
		decorator := a.State.ActDecorator[count-1]
		decorator(a.State)
	}
	a.State.CenterRowInit()
	for i := 0; i < 5; i++ {
		a.State.Draw()
	}
	fmt.Println("cardinhand", len(a.State.CardsInHand))
}

// create customized DefaultGamestate using several parameters
// starterDeckSets and centerDecksets refers to string const in factory.
// decorators is from decorator package
func CustomizedDefaultGamestate(starterDeckSets, centerDeckSets []string, decorators []decorator.AbstractDecorator) cards.AbstractGamestate {
	defGameState := NewDefaultGamestate().(*DefaultGamestate)
	defGameState.cardPiker = nil
	defGameState.currentResource = cards.NewResource()
	defGameState.CardsPlayed = []cards.Card{}
	defGameState.TopicsListeners = map[string]*DummyEventListener{}
	defGameState.CenterCards = []cards.Card{}
	defGameState.CardsInHand = []cards.Card{}
	// d.cardPiker = &TextCardPicker{}
	defGameState.HitPoint = 20
	defGameState.CardsDiscarded = cards.Deck{}
	defGameState.CardsInCenterDeck = cards.Deck{}
	centerDeck := []cards.Card{}
	for _, setName := range centerDeckSets {
		newCards := factory.CardFactory(setName, defGameState)
		centerDeck = append(centerDeck, newCards...)
	}
	defGameState.CardsInCenterDeck.SetList(centerDeck)
	defGameState.CardsInDeck = cards.Deck{}
	starterDeck := []cards.Card{}
	for _, setName := range starterDeckSets {
		newCards := factory.CardFactory(setName, defGameState)
		starterDeck = append(starterDeck, newCards...)
	}
	defGameState.CardsInDeck.SetList(starterDeck)
	defGameState.CardsInDeck.Shuffle()

	var j cards.AbstractGamestate
	// j := decorators[0](&defGameState)
	for _, k := range decorators {
		j = k(defGameState)
	}
	if j == nil {
		return defGameState
	} else {
		return j
	}
}
