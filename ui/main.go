package main

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/decorator"
	"github/kharism/GuildSim_go/internal/factory"
	"github/kharism/GuildSim_go/internal/gamestate"
	"log"
	"os"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Game struct {
	musicPlayer   *Player
	musicPlayerCh chan *Player
}

var (
	err        error
	background *ebiten.Image
)

const (
	HAND_SCALE = 0.25

	ORI_CARD_WIDTH  = 450
	ORI_CARD_HEIGHT = 600

	SPRITE_WIDTH  = 450 * 0.25
	SPRITE_HEIGHT = 600 * 0.25

	MAIN_DECK_X = 30
	MAIN_DECK_Y = 600*0.75 - 30

	HAND_DIST_X     = 450*HAND_SCALE + 10
	HAND_START_X    = MAIN_DECK_X + HAND_DIST_X
	HAND_START_Y    = MAIN_DECK_Y
	STATE_MAIN_MENU = "mainmenu"
	STATE_MAIN_GAME = "maingame"
	dpi             = 72

	PLAYED_START_X = 30
	PLAYED_START_Y = 600*3/4 - 200

	CENTER_DECK_START_X = 30
	CENTER_DECK_START_Y = 75

	ITEM_ICON_START_X = 20
	ITEM_ICON_START_Y = 0

	QUEST_ICON_START_X = 650
	QUEST_ICON_START_Y = -40

	DISCARD_NA_SOURCE_X = 350 //600 - 450*3/4
	DISCARD_NA_SOURCE_Y = 150 //300 - 300*3/4

	CENTER_START_X = CENTER_DECK_START_X + HAND_DIST_X
	CENTER_START_Y = CENTER_DECK_START_Y

	BANISHED_START_X = 1100
	BANISHED_START_Y = CENTER_START_Y

	DISCARD_START_X = 1000
	DISCARD_START_Y = MAIN_DECK_Y

	ENDTURN_START_X = 1100
	ENDTURN_START_Y = MAIN_DECK_Y

	DMG_START_X = 500
	DMG_START_Y = 600*0.75 + 50
)

var mainMenu AbstractEbitenState
var mainGame AbstractEbitenState
var currentState AbstractEbitenState
var mplusNormalFont font.Face
var mplusResource font.Face
var mplusDamage font.Face
var tooltipText font.Face
var audioContext *audio.Context

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	tooltipText, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    0.001,
		DPI:     10,
		Hinting: font.HintingFull,
	})
	mplusResource, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	mplusDamage, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

type exitAction struct{}

func (e *exitAction) DoAction() {
	//fmt.Println("Game over")
	//os.Exit(0)
	ll := mainGame.(*MainGameState)
	ll.mutex.Lock()
	ll.currentSubState = ll.gameoverState
	ll.cardsPlayed = []*EbitenCard{}
	ll.cardInHand = []*EbitenCard{}
	ll.mutex.Unlock()
}
func AttachGameOverListener(state cards.AbstractGamestate) cards.AbstractGamestate {
	quit := exitAction{}
	gameoverlistener := cards.NewStillAliveListener(state, &quit)
	onTakeDamage := &onTakeDamage{mainGameState: mainGame.(*MainGameState)}
	state.AttachListener(cards.EVENT_TAKE_DAMAGE, gameoverlistener)
	state.AttachListener(cards.EVENT_TAKE_DAMAGE, onTakeDamage)
	state.AttachListener(cards.EVENT_HEAL_DAMAGE, onTakeDamage)
	return state
}

func AttachDrawMainDeckListener(state cards.AbstractGamestate) cards.AbstractGamestate {
	draw := &OnDrawAction{mainGameState: mainGame.(*MainGameState)}
	state.AttachListener(cards.EVENT_CARD_DRAWN, draw)
	// fmt.Println("Done Attaching")
	return state
}
func AttachCardPlayedListener(state cards.AbstractGamestate) cards.AbstractGamestate {
	onPlayAction := &OnPlayAction{mainGameState: mainGame.(*MainGameState)}
	state.AttachListener(cards.EVENT_ATTR_CARD_PLAYED, onPlayAction)
	return state
}
func AttachCardDiscardListener(state cards.AbstractGamestate) cards.AbstractGamestate {
	onDiscardAction := &onDiscardAction{mainGameState: mainGame.(*MainGameState)}
	onBanishAction := &onBanishAction{mainGameState: mainGame.(*MainGameState)}
	state.AttachListener(cards.EVENT_ATTR_CARD_DISCARDED, onDiscardAction)
	state.AttachListener(cards.EVENT_ATTR_CARD_BANISHED, onBanishAction)
	return state
}
func AttachReturnToCenterListener(state cards.AbstractGamestate) cards.AbstractGamestate {
	onReturn := &onGotoCenterDeckAction{mainGameState: mainGame.(*MainGameState)}
	state.AttachListener(cards.EVENT_CARD_GOTO_CENTER, onReturn)
	return state
}
func AttachCenterCardDrawnListener(state cards.AbstractGamestate) cards.AbstractGamestate {
	onCenterDrawn := &onCenterDrawAction{mainGameState: mainGame.(*MainGameState)}
	state.AttachListener(cards.EVENT_CARD_DRAWN_CENTER, onCenterDrawn)
	return state
}
func AttachCenterCardRecDefExp(state cards.AbstractGamestate) cards.AbstractGamestate {
	onExplore := &onExplorationAction{mainGameState: mainGame.(*MainGameState)}
	onDefeat := &onDefeatAction{mainGameState: mainGame.(*MainGameState)}
	onRecruit := &onRecruitAction{mainGameState: mainGame.(*MainGameState)}
	onDisarm := &onDisarmAction{mainGameState: mainGame.(*MainGameState)}
	onItemAdd := &onItemAdd{mainGameState: mainGame.(*MainGameState)}
	onCardStacked := &onCardStacked{mainGameState: mainGame.(*MainGameState)}
	onPrePunish := &onPrePunish{mainGameState: mainGame.(*MainGameState)}
	onAddResource := &onChangeResource{mainGameState: mainGame.(*MainGameState)}
	ff := &onLimiterAttach{mainGameState: mainGame.(*MainGameState)}
	fg := &onLimiterDetach{mainGameState: mainGame.(*MainGameState)}
	detachAct := &onDetachAction{mainGameState: mainGame.(*MainGameState)}
	preTrap := &onBeforeTrap{MainGameState: mainGame.(*MainGameState)}
	onBossDefeated := &onBossDefeated{mainGameState: mainGame.(*MainGameState), bossDefeatedAction: gamestate.BossDefeatedAction{State: state.(*gamestate.DefaultGamestate)}}
	state.AttachListener(cards.EVENT_CARD_EXPLORED, onExplore)
	state.AttachListener(cards.EVENT_CARD_RECRUITED, onRecruit)
	state.AttachListener(cards.EVENT_CARD_DEFEATED, onDefeat)
	state.AttachListener(cards.EVENT_TRAP_REMOVED, onDisarm)
	state.AttachListener(cards.EVENT_ITEM_ADDED, onItemAdd)
	state.AttachListener(cards.EVENT_CARD_STACKED, onCardStacked)
	state.AttachListener(cards.EVENT_ATTACH_LIMITER, ff)
	state.AttachListener(cards.EVENT_DETACH_LIMITER, fg)
	state.AttachListener(cards.EVENT_REMOVE_OVERLAY, detachAct)
	state.AttachListener(cards.EVENT_BEFORE_PUNISH, onPrePunish)
	state.AttachListener(cards.EVENT_ADD_RESOURCE, onAddResource)
	state.AttachListener(cards.EVENT_BOSS_DEFEATED, onBossDefeated)
	state.AttachListener(cards.EVENT_BEFORE_TRAP, preTrap)
	return state
}
func (g *Game) ChangeState(stateName string) {
	switch stateName {
	case STATE_MAIN_GAME:
		// play different bgsound
		player, err := NewPlayer(audioContext, "sound/the-celtic-handmaiden-147078.mp3", typeMP3)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}

		g.musicPlayer.audioPlayer.Pause()
		g.musicPlayer = player
		cards.ResetPool()
		starterDeckSet := []string{factory.SET_STARTER_DECK}
		centerDeckSet := []string{factory.SET_CENTER_DECK_1}
		decorators := []decorator.AbstractDecorator{decorator.AttachTombOfForgottenMonarch, decorator.AttachTreasure, decorator.AttachProgressionCounter,
			AttachGameOverListener, AttachDrawMainDeckListener, AttachCardPlayedListener,
			AttachCardDiscardListener, AttachCenterCardDrawnListener, AttachCenterCardRecDefExp,
			AttachReturnToCenterListener,
		}
		defaultGamestate := gamestate.CustomizedDefaultGamestate(starterDeckSet, centerDeckSet, decorators)
		defaultGamestate.AddActDecorator(decorator.AttachHuntForDragonLord)
		// mainGame = NewMainGameState(g)
		mm := mainGame.(*MainGameState)
		mm.defaultGamestate = defaultGamestate.(*gamestate.DefaultGamestate)
		mm.hp = defaultGamestate.GetCurrentHP()
		mm.defaultGamestate.SetCardPicker(mm.cardPicker)
		mm.defaultGamestate.SetDetailViewer(mm.detailState)
		mm.defaultGamestate.SetBoolPicker(mm.boolPicker)
		// mm.defaultGamestate.AddResource(cards.RESOURCE_NAME_COMBAT, 20)
		mm.mainState = &mainMainState{m: mm, mutex: &sync.Mutex{}}
		mm.currentSubState = mm.mainState
		mm.mainState.Reset()
		// uncursePotion := item.NewBanishPotion(mm.defaultGamestate)
		// mm.defaultGamestate.AddItem(&uncursePotion)
		// stunCurse := cards.NewShockCurse(mm.defaultGamestate)
		// mm.defaultGamestate.CardsDiscarded.Push(&stunCurse)

		// slowTrap := cards.NewDragonValley(mm.defaultGamestate)
		// elephantDjinn := cards.NewAvalanceDragon(mm.defaultGamestate)
		// mm.defaultGamestate.AddCardToCenterDeck(cards.DISCARD_SOURCE_NAN, false, &elephantDjinn)
		// mm.defaultGamestate.AddResource(cards.RESOURCE_NAME_EXPLORATION, 12)
		// mm.defaultGamestate.AddResource(cards.RESOURCE_NAME_COMBAT, 12)
		// jj := cards.NewAggroDjinn(mm.defaultGamestate)
		// mm.defaultGamestate.StackCards(cards.DISCARD_SOURCE_NAN, &jj)

		// jj := cards.NewForgottenMonarchP2(mm.defaultGamestate)
		// mm.defaultGamestate.CardsInCenterDeck.Stack(&jj)
		// mm.defaultGamestate.TakeDamage(40)
		// wl := cards.NewWingedLion(mm.defaultGamestate)
		// dw := cards.NewDeadweight(mm.defaultGamestate)
		// kk := cards.NewRookieMage(mm.defaultGamestate)
		// slimeRoom := cards.NewSlimeRoom(mm.defaultGamestate)
		// boulder := cards.NewDragonLord(mm.defaultGamestate)
		// spikeFloor := cards.NewSpikeFloor(mm.defaultGamestate)
		// lair := cards.NewGoblinSmallLairArea(mm.defaultGamestate)
		// heal := item.NewHealingPotion(defaultGamestate)
		// ll := append(mm.defaultGamestate.CardsInCenterDeck.List()[:3], &spikeFloor)
		// rest := mm.defaultGamestate.CardsInCenterDeck.List()
		// iceWyvern := cards.NewIceWyvern(mm.defaultGamestate)
		// mm.defaultGamestate.CardsInCenterDeck.SetList(append([]cards.Card{&iceWyvern}, rest...))
		// mm.defaultGamestate.ItemCards = append(mm.defaultGamestate.ItemCards, &heal)
		// mm.defaultGamestate.CardsInHand = append(mm.defaultGamestate.CardsInHand, &spikeFloor)
		// list := mm.defaultGamestate.CardsInDeck.List()
		// newList := append(list[0:3], &kk)
		// newList = append(newList, list[4:]...)
		// mm.defaultGamestate.CardsInDeck.SetList(newList)
		// mm.defaultGamestate.CardsInDeck.Stack(&kk)
		// newDeck := []cards.Card{&boulder}
		// rookieCard := NewEbitenCardFromCard(&spikeFloor)
		// rookieCard.x = HAND_START_X
		// rookieCard.y = HAND_START_Y
		// mm.cardInHand = append(mm.cardInHand, rookieCard)
		currentState = mainGame
		// vampire := cards.NewRagingVampire(mm.defaultGamestate)
		// mm.defaultGamestate.CardsInCenterDeck.Stack(vampire)
		mm.defaultGamestate.CenterRowInit()
		// for i := 0; i < 10; i++ {
		// 	hh := cards.NewArcher(mm.defaultGamestate)
		// 	mm.defaultGamestate.CardsDiscarded.Stack(&hh)
		// 	hh2 := cards.NewAggroDjinn(mm.defaultGamestate)
		// 	mm.defaultGamestate.CardsDiscarded.Stack(&hh2)
		// }
		// mm.defaultGamestate.CardsInCenterDeck.Stack(&boulder)
		// mm.defaultGamestate.CardsInCenterDeck.Stack(&spikeFloor)
		// mm.defaultGamestate.CardsInCenterDeck.Stack(&slimeRoom)

		// rookieMage := cards.NewRookieMage(mm.defaultGamestate)
		// stunCurse1 := cards.NewStunCurse(mm.defaultGamestate)
		// stunCurse2 := cards.NewStunCurse(mm.defaultGamestate)
		// stunCurse3 := cards.NewStunCurse(mm.defaultGamestate)
		// mm.defaultGamestate.GetMainDeck().Stack(vampire)
		// mm.defaultGamestate.GetMainDeck().Stack(&rookieMage)
		// mm.defaultGamestate.GetMainDeck().Stack(&stunCurse3)
		mm.defaultGamestate.BeginTurn()
	case STATE_MAIN_MENU:
		player, err := NewPlayer(audioContext, "sound/perfect-beauty-191271.mp3", typeMP3)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}

		g.musicPlayer.audioPlayer.Pause()
		g.musicPlayer = player
		currentState = mainMenu
	}
}
func (g *Game) Update() error {
	g.musicPlayer.update()
	return currentState.Update()
	// return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	currentState.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1200, 600
}

var game *Game

func main() {
	audioContext = audio.NewContext(sampleRate)
	ebiten.SetWindowSize(1200, 600)
	ebiten.SetWindowTitle("Hello, World!")
	game = &Game{
		musicPlayerCh: make(chan *Player),
	}
	m, err := NewPlayer(audioContext, "./sound/perfect-beauty-191271.mp3", typeMP3)
	if err != nil {
		fmt.Println(err.Error())
		return
		// return nil, err
	}
	game.musicPlayer = m
	mainMenu = NewMainMenuState(game)
	mainGame = NewMainGameState(game)

	currentState = mainMenu
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
