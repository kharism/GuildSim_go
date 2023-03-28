package cards

import "github/kharism/GuildSim_go/internal/observer"

// type of events
var EVENT_CARD_PLAYED = "CardPlayed"
var EVENT_CARD_DRAWN = "CardDraw"
var EVENT_CARD_DRAWN_CENTER = "CardDrawCenter"
var EVENT_CARD_GOTO_CENTER = "CardGotoCenter"
var EVENT_CARD_EXPLORED = "CardExplored"
var EVENT_CARD_DEFEATED = "CardDefeated"
var EVENT_CARD_RECRUITED = "CardRecruited"
var EVENT_CARD_DISCARDED = "CardDiscarded"
var EVENT_CARD_BANISHED = "CardBanished"
var EVENT_TAKE_DAMAGE = "TakeDamage"
var EVENT_HEAL_DAMAGE = "HealDamage"

var EVENT_START_OF_TURN = "BeginTurn"
var EVENT_END_OF_TURN = "EndTurn"

// this is the key of map[string]interface
var EVENT_ATTR_CARD_PLAYED = "CardPlayed"
var EVENT_ATTR_CARD_DRAWN = "CardDrawn"
var EVENT_ATTR_CARD_EXPLORED = "CardExplored"
var EVENT_ATTR_CARD_DEFEATED = "CardDefeated"
var EVENT_ATTR_CARD_RECRUITED = "CardRecruited"
var EVENT_ATTR_CARD_DISCARDED = "CardDiscarded"
var EVENT_ATTR_CARD_BANISHED = "CardBanished"
var EVENT_ATTR_CARD_GOTO_CENTER = "CardGotoCenter"
var EVENT_ATTR_CARD_TAKE_DAMAGE = "CardTakeDamage"
var EVENT_ATTR_CARD_TAKE_DAMAGE_AMMOUNT = "CardTakeDamageAmt"
var EVENT_ATTR_CARD_HEAL_DAMAGE = "CardHealDamage"

var EVENT_ATTR_DAMAGE_AMMT = "DamageAmount"

var EVENT_ATTR_DISCARD_SOURCE = "DiscardSource"

const (
	DISCARD_SOURCE_HAND     = "hand"
	DISCARD_SOURCE_PLAYED   = "played"
	DISCARD_SOURCE_CENTER   = "center"
	DISCARD_SOURCE_COOLDOWN = "cooldown"
	DISCARD_SOURCE_NAN      = "nan"
	DISCARD_SOURCE_DISCARD  = "discard"
)

type AbstractGamestate interface {
	GetPlayedCards() []Card
	GetCardInHand() []Card
	GetCenterCard() []Card
	GetCooldownCard() []Card
	RecruitCard(c Card)
	DiscardCard(c Card, source string)
	BanishCard(c Card, source string)
	DefeatCard(c Card)
	// just play card from no particular location and added it to list of played card
	// It will assume the card is played from hand and try to remove cards from hand if possible
	// the card will not automatically go to discard/cooldown pile
	// otherwise remove the card accordingly
	PlayCard(c Card)
	Explore(c Card)

	// item related stuff
	ListItems() []Card
	RemoveItem(c Card)
	RemoveItemIndex(i int)
	ConsumeItem(c Consumable)

	// end turn, remove event listener attached by played cards, remove resources except money+reputation,
	// take punishment etc
	EndTurn()
	// begin turn, draw 5 and THEN apply any eff that that need to be applied at the start of the turn
	BeginTurn()

	// damage
	GetCurrentHP() int

	// take damage, the parameter can also take negative damage which means heals
	// it also trigger takeDamage or healDamage event
	TakeDamage(int)

	// remove cards
	RemoveCardFromHand(c Card)
	RemoveCardFromHandIdx(i int)
	RemoveCardFromCenterRow(c Card)
	RemoveCardFromCenterRowIdx(i int)
	RemoveCardFromCooldown(c Card)
	RemoveCardFromCooldownIdx(i int)

	// return a card drawn from central deck
	ReplaceCenterCard() Card
	// init center row
	CenterRowInit()
	AppendCenterCard(c Card)

	// its pot of greed but halved
	Draw()

	// get abstract card picker
	GetCardPicker() AbstractCardPicker
	SetCardPicker(AbstractCardPicker)

	GetBoolPicker() AbstractBoolPicker
	SetBoolPicker(AbstractBoolPicker)

	// make center deck thicker
	AddCardToCenterDeck(source string, shuffle bool, c ...Card)

	AttachListener(eventName string, l observer.Listener)
	RemoveListener(eventName string, l observer.Listener)
	GetCurrentResource() Resource
	AddResource(name string, amount int)
	PayResource(cost Cost)
}
