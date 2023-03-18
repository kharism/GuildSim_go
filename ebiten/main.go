package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

var (
	err        error
	background *ebiten.Image
)

const (
	HAND_SCALE = 0.25

	SPRITE_WIDTH  = 450 * 0.25
	SPRITE_HEIGHT = 600 * 0.25

	MAIN_DECK_X = 30
	MAIN_DECK_Y = 600*0.75 - 30

	HAND_DIST_X     = 450*HAND_SCALE + 30
	HAND_START_X    = MAIN_DECK_X + HAND_DIST_X
	HAND_START_Y    = MAIN_DECK_Y
	STATE_MAIN_MENU = "mainmenu"
	STATE_MAIN_GAME = "maingame"
)

var mainMenu AbstractEbitenState
var mainGame AbstractEbitenState
var currentState AbstractEbitenState

func init() {

}
func (g *Game) ChangeState(stateName string) {
	switch stateName {
	case STATE_MAIN_GAME:
		currentState = mainGame
	case STATE_MAIN_MENU:
		currentState = mainMenu
	}
}
func (g *Game) Update() error {
	return currentState.Update()
	// return nil
}

var ShowCardDetail = false

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
