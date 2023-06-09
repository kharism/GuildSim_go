package main

import (
	"github/kharism/GuildSim_go/internal/cards"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	csg "github.com/kharism/golang-csg/core"
)

type EbitenCard struct {
	image     *ebiten.Image
	card      cards.Card
	oriWidth  int
	oriHeight int
	// current position
	x float64
	y float64
	// velocity of card movement
	vx float64
	vy float64
	// target position if card moved
	tx float64
	ty float64
	// translation on x axis due to dragging
	x_drag int

	// syncinc stuff
	mutex *sync.Mutex

	// animation stuff
	CurrMove       *MoveAnimation
	AnimationQueue []*MoveAnimation
}

func (e *EbitenCard) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(0.25, 0.25)
	op.GeoM.Translate(0, 0)
	// op.GeoM.Translate(MAIN_DECK_X, MAIN_DECK_Y)
	op.GeoM.Translate(float64(e.x), float64(e.y))

	op.GeoM.Translate(float64(e.x_drag), 0)
	// op.GeoM.Translate(float64(e.x), float64(e.y))
	screen.DrawImage(e.image, op)
}
func (e *EbitenCard) AddAnimation(animation ...*MoveAnimation) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.AnimationQueue = append(e.AnimationQueue, animation...)
}
func (e *EbitenCard) Update() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.CurrMove == nil && len(e.AnimationQueue) > 0 {
		e.CurrMove = e.AnimationQueue[0]
		e.AnimationQueue = e.AnimationQueue[1:]
		if e.CurrMove.SleepPre != 0 {
			time.Sleep(e.CurrMove.SleepPre)
		}
		e.tx = e.CurrMove.tx
		e.ty = e.CurrMove.ty
		vx := float64(e.tx - e.x)
		vy := float64(e.ty - e.y)
		speedVector := csg.NewVector(vx, vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
		e.vx = speedVector.X
		e.vy = speedVector.Y
	}
	e.x += e.vx
	e.y += e.vy
	// fmt.Println(e.x, e.y)
	if math.Abs(float64(e.tx-e.x))+math.Abs(float64(e.ty-e.y)) < 15 {
		if e.CurrMove != nil && e.CurrMove.DoneFunc != nil {
			if e.CurrMove.SleepPost != 0 {
				//time.Sleep(e.CurrMove.SleepPost)
			}
			e.CurrMove.DoneFunc()
		}
		if len(e.AnimationQueue) == 0 {
			e.x = e.tx
			e.y = e.ty
			e.vx = 0
			e.vy = 0
			e.CurrMove = nil
		} else {
			e.CurrMove = e.AnimationQueue[0]
			e.AnimationQueue = e.AnimationQueue[1:]
			if e.CurrMove.SleepPre != 0 {
				//time.Sleep(e.CurrMove.SleepPre)
			}
			e.tx = e.CurrMove.tx
			e.ty = e.CurrMove.ty
			vx := float64(e.tx - e.x)
			vy := float64(e.ty - e.y)
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
			e.vx = speedVector.X
			e.vy = speedVector.Y
		}

	}
}
