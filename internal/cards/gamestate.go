package cards

import "github/kharism/GuildSim_go/internal/observer"

// type of events
var EVENT_CARD_PLAYED = "CardPlayed"
var EVENT_CARD_EXPLORED = "CardExplored"
var EVENT_CARD_DEFEATED = "CardDefeated"
var EVENT_CARD_BANISHED = "CardBanished"
var EVENT_TAKE_DAMAGE = "TakeDamage"

// this is the key of map[string]interface
var EVENT_ATTR_CARD_PLAYED = "CardPlayed"
var EVENT_ATTR_CARD_EXPLORED = "CardExplored"
var EVENT_ATTR_CARD_DEFEATED = "CardDefeated"
var EVENT_ATTR_CARD_BANISHED = "CardBanished"

var EVENT_ATTR_DAMAGE_AMMT = "DamageAmount"

type AbstractGamestate interface {
	GetPlayedCards() []Card
	GetCardInHand() []Card
	GetCenterCard() []Card
	RecruitCard(c Card)
	DiscardCard(c Card)
	BanishCard(c Card)
	DefeatCard(c Card)
	// just play card from no particular location and added it to list of played card
	// It will assume the card is played from hand and try to remove cards from hand if possible
	// the card will not automatically go to discard/cooldown pile
	// otherwise remove the card accordingly
	PlayCard(c Card)
	Explore(c Card)

	// end turn, remove event listener attached by played cards, remove resources except money+reputation,
	// take punishment etc
	EndTurn()

	// damage
	GetCurrentHP() int
	TakeDamage(int)

	// remove cards
	RemoveCardFromHand(c Card)
	RemoveCardFromHandIdx(i int)
	RemoveCardFromCenterRow(c Card)
	RemoveCardFromCenterRowIdx(i int)

	// return a card drawn from central deck
	ReplaceCenterCard() Card
	// init center row
	CenterRowInit()

	// its pot of greed but halved
	Draw()

	// get abstract card picker
	GetCardPicker() AbstractCardPicker

	// make center deck thicker
	AddCardToCenterDeck(c ...Card)

	AttachListener(eventName string, l observer.Listener)
	RemoveListener(eventName string, l observer.Listener)
	GetCurrentResource() Resource
	AddResource(name string, amount int)
	PayResource(cost Cost)
}
