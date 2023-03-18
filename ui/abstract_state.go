package main

import "github.com/hajimehoshi/ebiten/v2"

//
type AbstractEbitenState interface {
	// execute this method on update
	Update() error
	Draw(screen *ebiten.Image)
}
type AbstractStateChanger interface {
	ChangeState(newState string)
}
