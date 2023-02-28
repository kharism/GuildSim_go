package cards

import "github/kharism/GuildSim_go/internal/observer"

// type of events
var EVENT_CARD_PLAYED = "CardPlayed"
var EVENT_CARD_EXPLORED = "CardExplored"
var EVENT_CARD_DEFEATED = "CardDefeated"
var EVENT_CARD_BANISHED = "CardBanished"

// this is the key of map[string]interface
var EVENT_ATTR_CARD_PLAYED = "CardPlayed"
var EVENT_ATTR_CARD_EXPLORED = "CardExplored"
var EVENT_ATTR_CARD_DEFEATED = "CardDefeated"
var EVENT_ATTR_CARD_BANISHED = "CardBanished"

type AbstractGamestate interface {
	GetPlayedCards() []Card
	GetCardInHand() []Card
	GetCenterCard() []Card
	RecruitCard(c Card)
	DiscardCard(c Card)
	BanishCard(c Card)
	DefeatCard(c Card)
	PlayCard(c Card)
	Explore(c Card)

	// return a card drawn from central deck
	ReplaceCenterCard() Card
	// init center row
	CenterRowInit()

	// its pot of greed but halved
	Draw()

	AttachListener(eventName string, l observer.Listener)
	RemoveListener(eventName string, l observer.Listener)
	GetCurrentResource() Resource
	AddResource(name string, amount int)
	PayResource(cost Cost)
}
