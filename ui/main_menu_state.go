package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MainMenuState struct {
	bgImage      *ebiten.Image
	newGameBtn   *ebiten.Image
	exitBtn      *ebiten.Image
	stateChanger AbstractStateChanger
}

func NewMainMenuState(stateChanger AbstractStateChanger) AbstractEbitenState {
	bgImage, _, err := ebitenutil.NewImageFromFile("img/menu_screen.png")
	newGameBtn, _, err := ebitenutil.NewImageFromFile("img/new_game.png")
	exitBtn, _, err := ebitenutil.NewImageFromFile("img/exit.png")
	if err != nil {
		log.Fatal(err)
	}
	return &MainMenuState{bgImage: bgImage, newGameBtn: newGameBtn, exitBtn: exitBtn, stateChanger: stateChanger}
}
func (m *MainMenuState) Update() error {
	return nil
}

const (
	buttonWidth  = 150
	buttonHeight = 75
)

func (m *MainMenuState) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(m.bgImage, op)
	op2 := &ebiten.DrawImageOptions{}
	newGamePosX, newGamePosY := 600-100, 300-80
	op2.GeoM.Translate(float64(newGamePosX), float64(newGamePosY))

	screen.DrawImage(m.newGameBtn, op2)
	op3 := &ebiten.DrawImageOptions{}
	exitPosX, exitPosY := 600-100, 300
	op3.GeoM.Translate(float64(exitPosX), float64(exitPosY))

	screen.DrawImage(m.exitBtn, op3)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > newGamePosX && x < newGamePosX+buttonWidth {
			if y > newGamePosY && y < newGamePosY+buttonHeight {
				fmt.Println("NewGame")
				m.stateChanger.ChangeState(STATE_MAIN_GAME)
			} else if y > exitPosY && y < exitPosY+buttonHeight {
				// fmt.Println("Exit")
				os.Exit(0)
			}
		}
	}
}
