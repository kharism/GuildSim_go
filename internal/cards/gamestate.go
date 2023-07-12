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
var EVENT_CARD_STACKED = "CardStacked"
var EVENT_TAKE_DAMAGE = "TakeDamage"
var EVENT_HEAL_DAMAGE = "HealDamage"
var EVENT_ITEM_ADDED = "ItemAdded"
var EVENT_ITEM_REMOVED = "ItemRemoved"
var EVENT_TRAP_REMOVED = "TrapRemoved"
var EVENT_ATTACH_LIMITER = "AttachLimiter"
var EVENT_DETACH_LIMITER = "DetachLimiter"
var EVENT_BEFORE_PUNISH = "BeforePunish"
var EVENT_ADD_RESOURCE = "AddResource"

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
var EVENT_ATTR_CARD_STACKED = "CardStacked"
var EVENT_ATTR_CARD_GOTO_CENTER = "CardGotoCenter"
var EVENT_ATTR_CARD_TAKE_DAMAGE = "CardTakeDamage"
var EVENT_ATTR_CARD_TAKE_DAMAGE_AMMOUNT = "CardTakeDamageAmt"
var EVENT_ATTR_CARD_HEAL_DAMAGE = "CardHealDamage"
var EVENT_ATTR_ITEM_ADDED = "ItemAdded"
var EVENT_ATTR_ITEM_REMOVED = "ItemRemoved"
var EVENT_ATTR_TRAP_REMOVED = "TrapRemoved"
var EVENT_ATTR_BEFORE_PUNISH_CARD = "CardBeforePunish"
var EVENT_ATTR_ADD_RESOURCE_NAME = "AttrAddResourceName"
var EVENT_ATTR_ADD_RESOURCE_AMOUNT = "AttrAddResourceAmount"

var EVENT_ATTR_LIMITER = "Limiter"
var EVENT_ATTR_LIMITER_ACTION = "LimiterAction"

var EVENT_ATTR_DAMAGE_AMMT = "DamageAmount"

var EVENT_ATTR_DISCARD_SOURCE = "DiscardSource"

var EVENT_MINIBOSS_DEFEATED = "MinibossDefeated"
var EVENT_BOSS_DEFEATED = "BossDefeated"

var EVENT_ATTR_BOSS_DEFEATED_COUNT = "BossDefeatedCount"

const (
	ACTION_DRAW    = "Draw"
	ACTION_DEFEAT  = "Defeat"
	ACTION_EXPLORE = "Explore"
	ACTION_RECRUIT = "Recruit"
	ACTION_DISARM  = "Disarm"
)

const (
	DISCARD_SOURCE_HAND        = "hand"
	DISCARD_SOURCE_PLAYED      = "played"
	DISCARD_SOURCE_CENTER      = "center"
	DISCARD_SOURCE_CENTER_DECK = "centerDeck"
	DISCARD_SOURCE_COOLDOWN    = "cooldown"
	DISCARD_SOURCE_NAN         = "nan"
	DISCARD_SOURCE_DISCARD     = "discard"
)
const (
	RARITY_COMMON = 0b001
	RARITY_RARE   = 0b010
	RARITY_SRARE  = 0b100
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

	// thsi is method to call when remove trap in center row.
	Disarm(c Card)
	// just play card from no particular location and added it to list of played card
	// It will assume the card is played from hand and try to remove cards from hand if possible
	// the card will not automatically go to discard/cooldown pile
	// otherwise remove the card accordingly
	PlayCard(c Card)
	Explore(c Card)

	MutexLock()
	MutexUnlock()

	// rewardGeneratorFunctions
	GenerateRandomPotion(rarity int) Card
	GenerateRandomRelic(rarity int) Card

	// item related stuff
	ListItems() []Card
	RemoveItem(c Card)
	RemoveItemIndex(i int)
	AddItem(c Card)
	ConsumeItem(c Consumable)

	// end turn, remove event listener attached by played cards, remove resources except money+reputation,
	// take punishment etc
	EndTurn()
	// begin turn, draw 5 and THEN apply any eff that that need to be applied at the start of the turn
	BeginTurn()

	// damage
	GetCurrentHP() int

	// Take damage, the parameter can also take negative damage which means heals.
	// It also trigger takeDamage or healDamage event.
	// Since we added Block resource, deduce the damage to the block resource.
	// like if a monster call for TakeDamage(3) damage and you have 2 block, then
	// player HP will be reduced by 1 and their block will be 0
	TakeDamage(int)

	// remove cards
	RemoveCardFromHand(c Card)
	RemoveCardFromHandIdx(i int)
	RemoveCardFromCenterRow(c Card)
	RemoveCardFromCenterRowIdx(i int)
	UpdateCenterCard(c Card)
	RemoveCardFromCooldown(c Card)
	RemoveCardFromCooldownIdx(i int)

	// return a card drawn from central deck
	ReplaceCenterCard() Card
	// peek what card at the top of the center deck
	PeekCenterCard() Card
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

	SetDetailViewer(AbstractDetailViewer)
	GetDetailViewer() AbstractDetailViewer

	// make center deck thicker
	AddCardToCenterDeck(source string, shuffle bool, c ...Card)

	// put cards to top of main deck, the order matter
	// c[0] will be stacked first, then c[1], then c[2], etc
	StackCards(source string, c ...Card)
	ShuffleMainDeck()

	AttachListener(eventName string, l observer.Listener)
	RemoveListener(eventName string, l observer.Listener)
	NotifyListener(eventname string, data map[string]interface{})

	// legalchecks. TODO: implement this on DefGamestate and Dummy
	AttachLegalCheck(actionName string, lc LegalChecker)
	DetachLegalCheck(actionName string, lc LegalChecker)
	LegalCheck(actionName string, data interface{}) bool

	GetCurrentResource() Resource
	AddResource(name string, amount int)
	PayResource(cost Cost)

	// this function will decorate the current gamestate since we now have multi act gamenow
	ActDecorators() []func(AbstractGamestate) AbstractGamestate
	AddActDecorator(func(AbstractGamestate) AbstractGamestate)
}
