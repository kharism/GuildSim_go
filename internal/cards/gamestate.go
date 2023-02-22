package cards

import "github/kharism/GuildSim_go/internal/observer"

// type of events
var EVENT_CARD_PLAYED = "CardPlayed"

// this is the key of map[string]interface
var EVENT_ATTR_CARD_PLAYED = "CardPlayed"

type AbstractGamestate interface {
	GetPlayedCards() []Card
	GetCardInHand() []Card
	GetCenterCard() []Card
	RecruitCard(c Card)
	DiscardCard(c Card)
	BanishCard(c Card)
	DefeatCard(c Card)
	PlayCard(c Card)
	AttachListener(eventName string, l observer.Listener)
	RemoveListener(eventName string, l observer.Listener)
	GetCurrentResource() Resource
	AddResource(name string, amount int)
}
