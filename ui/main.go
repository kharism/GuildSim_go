package main

import (
	"github/kharism/GuildSim_go/internal/cards"
	"github/kharism/GuildSim_go/internal/decorator"
	"github/kharism/GuildSim_go/internal/factory"
	"github/kharism/GuildSim_go/internal/gamestate"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

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

	HAND_DIST_X     = 450*HAND_SCALE + 30
	HAND_START_X    = MAIN_DECK_X + HAND_DIST_X
	HAND_START_Y    = MAIN_DECK_Y
	STATE_MAIN_MENU = "mainmenu"
	STATE_MAIN_GAME = "maingame"
	dpi             = 72
)

var mainMenu AbstractEbitenState
var mainGame AbstractEbitenState
var currentState AbstractEbitenState
var mplusNormalFont font.Face
var mplusResource font.Face

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
	mplusResource, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

type exitAction struct{}

func (e *exitAction) DoAction() {
	//fmt.Println("Game over")
	//os.Exit(0)
	currentState = mainMenu
}
func AttachGameOverListener(state cards.AbstractGamestate) cards.AbstractGamestate {
	quit := exitAction{}
	gameoverlistener := cards.NewStillAliveListener(state, &quit)
	state.AttachListener(cards.EVENT_TAKE_DAMAGE, gameoverlistener)
	return state
}

func AttachDrawMainDeckListener(state cards.AbstractGamestate) cards.AbstractGamestate {
	draw := &OnDrawAction{mainGameState: mainGame.(*MainGameState)}
	state.AttachListener(cards.EVENT_CARD_DRAWN, draw)
	// fmt.Println("Done Attaching")
	return state
}
func (g *Game) ChangeState(stateName string) {
	switch stateName {
	case STATE_MAIN_GAME:
		starterDeckSet := []string{factory.SET_STARTER_DECK}
		centerDeckSet := []string{factory.SET_CENTER_DECK_1}
		decorators := []decorator.AbstractDecorator{AttachGameOverListener, AttachDrawMainDeckListener}
		defaultGamestate := gamestate.CustomizedDefaultGamestate(starterDeckSet, centerDeckSet, decorators)
		mm := mainGame.(*MainGameState)
		mm.defaultGamestate = defaultGamestate.(*gamestate.DefaultGamestate)
		mm.defaultGamestate.SetCardPicker(mm.cardPicker)
		currentState = mainGame
		mm.defaultGamestate.BeginTurn()
	case STATE_MAIN_MENU:
		currentState = mainMenu
	}
}
func (g *Game) Update() error {
	return currentState.Update()
	// return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	currentState.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1200, 600
}

func main() {
	ebiten.SetWindowSize(1200, 600)
	ebiten.SetWindowTitle("Hello, World!")
	game := &Game{}
	mainMenu = NewMainMenuState(game)
	mainGame = NewMainGameState(game)

	currentState = mainMenu
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
