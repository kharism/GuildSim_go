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

	MouseIn bool

	overlays []EbitenCard

	// syncinc stuff
	mutex *sync.Mutex

	// animation stuff
	CurrMove       *MoveAnimation
	AnimationQueue []*MoveAnimation
}

const OVERLAY_MARGIN = 15

func (e *EbitenCard) IsMouseIn(mouseX, mouseY float64) bool {
	if mouseX > e.x && mouseX < e.x+HAND_DIST_X && mouseY > e.y && mouseY < e.y+SPRITE_HEIGHT {
		return true
	}
	return false
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
	// if e.card.GetName() == "Treasure" {
	if _, ok := e.card.(cards.Overlay); ok {
		l := e.card.(cards.Overlay)
		if l.HasOverlayCard() {
			overlays := l.GetOverlay()
			for idx, val := range overlays {
				op.GeoM.Reset()
				op.GeoM.Scale(0.25, 0.25)
				overlayPosX := float64(e.x)
				overlayPosY := float64(e.y) + float64(idx+1)*OVERLAY_MARGIN
				op.GeoM.Translate(overlayPosX, overlayPosY)
				pp := NewEbitenCardFromCard(val)
				screen.DrawImage(pp.image, op)
			}
		}
	}
	// }

}
func (e *EbitenCard) DrawTooltip(screen *ebiten.Image) {
	// jj := e.card.GetKeywords()
	// if e.MouseIn && len(jj) > 0 {
	// 	//ebitenutil.DrawRect(screen,)
	// 	// fmt.Println(e.card.GetName())
	// 	ebitenutil.DrawRect(screen, e.x+SPRITE_WIDTH, e.y, 140, 60, color.White)
	// 	if e.Keywords == nil {
	// 		keywordText := &EbitenText{}
	// 		keywordText.face = tooltipText
	// 		keywordText.x = e.x + SPRITE_WIDTH
	// 		keywordText.y = e.y + 12
	// 		keywordText.color = color.RGBA{255, 0, 0, 255}
	// 		keywordText.text = jj[0]
	// 		e.Keywords = keywordText

	// 	}
	// 	e.Keywords.Draw(screen)
	// }
}
func (e *EbitenCard) AddAnimation(animation ...*MoveAnimation) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.AnimationQueue = append(e.AnimationQueue, animation...)
}
func (e *EbitenCard) ReplaceCurrentAnim(animation *MoveAnimation) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.CurrMove = animation
	e.tx = e.CurrMove.tx
	e.ty = e.CurrMove.ty
	vx := float64(e.tx - e.x)
	vy := float64(e.ty - e.y)
	if vx != 0 || vy != 0 {
		speedVector := csg.NewVector(vx, vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
		e.vx = speedVector.X
		e.vy = speedVector.Y
	} else {
		e.vx = 0
		e.vy = 0
	}
}
func (e *EbitenCard) Update() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.CurrMove == nil && len(e.AnimationQueue) > 0 {
		e.CurrMove = e.AnimationQueue[0]
		e.AnimationQueue = e.AnimationQueue[1:]
		// fmt.Println("animation queue", e.card.GetName(), e.CurrMove)
		if e.CurrMove.SleepPre != 0 {
			time.Sleep(e.CurrMove.SleepPre)
		}
		e.tx = e.CurrMove.tx
		e.ty = e.CurrMove.ty
		vx := float64(e.tx - e.x)
		vy := float64(e.ty - e.y)
		if vx != 0 || vy != 0 {
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
			e.vx = speedVector.X
			e.vy = speedVector.Y
		} else {
			e.vx = 0
			e.vy = 0
		}

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
			if vy != 0 || vx != 0 {
				speedVector := csg.NewVector(vx, vy, 0)
				speedVector = speedVector.Normalize().MultiplyScalar(e.CurrMove.Speed)
				e.vx = speedVector.X
				e.vy = speedVector.Y
			} else {
				e.vx = 0
				e.vy = 0
			}

		}

	}
}
