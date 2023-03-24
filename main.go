package main

import (
	"github/kharism/GuildSim_go/internal/decorator"
	"github/kharism/GuildSim_go/internal/factory"
	"github/kharism/GuildSim_go/internal/gamestate"
	"github/kharism/GuildSim_go/internal/ui/text"
)

func main() {
	// build defaultgamestate here, let user parameterize this later
	starterDeckSet := []string{factory.SET_STARTER_DECK}
	centerDeckSet := []string{factory.SET_CENTER_DECK_1}
	decorators := []decorator.AbstractDecorator{decorator.AttachTombOfForgottenMonarch}
	defaultGamestate := gamestate.CustomizedDefaultGamestate(starterDeckSet, centerDeckSet, decorators)
	textUI := text.NewTextUIGamestate(defaultGamestate)
	defaultGamestate.SetCardPicker(textUI.GetCardPicker())
	// playerDeck := []cards.Card{}
	// for i := 0; i < 5; i++ {
	// 	kk := cards.NewRookieAdventurer(gamestate)
	// 	playerDeck = append(playerDeck, &kk)
	// }
	//gamestate.(*text.TextUIGamestate).Run()
	textUI.Run()
}
