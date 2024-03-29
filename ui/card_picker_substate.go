package main

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	CARDPICKER_START_X = 160
	CARDPICKER_START_Y = 40

	CARDPICKER_DIST_X = 120
	CARDPICKER_DIST_Y = 40
)

type cardPickState struct {
	m             *MainGameState
	cards         []cards.Card
	selectedCard  *EbitenCard
	optional      bool
	selectedIndex int
	pickedCards   chan (int)
	message       string
}
type cardListState struct {
	m     *MainGameState
	cards []cards.Card
}

func (c *cardPickState) Update() error {
	return nil
}
func (c *cardPickState) PickCard(list []cards.Card, message string) int {
	c.cards = list
	c.message = message
	// fmt.Println("Tunggu hasil")
	c.optional = false
	c.m.mutex.Lock()
	c.m.currentSubState = c
	for i := range c.m.cardsInCenter {
		fmt.Println(c.m.cardsInCenter[i].card.GetName(), c.m.cardsInCenter[i].x, c.m.cardsInCenter[i].y, c.m.cardsInCenter[i].tx, c.m.cardsInCenter[i].ty)
	}
	c.m.mutex.Unlock()
	pickedCards := <-c.pickedCards
	for i := range c.m.cardsInCenter {
		fmt.Println(c.m.cardsInCenter[i].card.GetName(), c.m.cardsInCenter[i].x, c.m.cardsInCenter[i].y, c.m.cardsInCenter[i].tx, c.m.cardsInCenter[i].ty)
	}
	// fmt.Println("Dapat hasil", pickedCards)
	return pickedCards
}
func (c *cardPickState) PickCardOptional(list []cards.Card, message string) int {
	c.cards = list
	c.message = message
	c.optional = true
	c.m.mutex.Lock()
	c.m.currentSubState = c
	c.m.mutex.Unlock()
	// fmt.Println("Tunggu hasil")
	pickedCards := <-c.pickedCards
	// fmt.Println("Dapat hasil", pickedCards)
	return pickedCards
}

func (c *cardPickState) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(c.m.bgImage2, op)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(1.3, 1.0)
	op2.GeoM.Translate(120, 0)
	screen.DrawImage(c.m.paperBg, op2)
	op3 := &ebiten.DrawImageOptions{}
	text.Draw(screen, c.message, mplusNormalFont, CARDPICKER_START_X, 40, color.RGBA{255, 255, 255, 255})
	colPerRow := 7
	cardList := []*EbitenCard{}
	for idx, cc := range c.cards {
		ebitenCard := NewEbitenCardFromCard(cc)
		op3.GeoM.Reset()
		op3.GeoM.Scale(HAND_SCALE, HAND_SCALE)
		col := (idx % colPerRow)
		row := (idx / colPerRow)
		// fmt.Println(row, col)
		ebitenCard.x = float64(CARDPICKER_START_X + CARDPICKER_DIST_X*col)
		ebitenCard.y = float64(CARDPICKER_START_Y + CARDPICKER_DIST_Y*row)
		// op3.GeoM.Translate(float64(CARDPICKER_START_X*col), float64(CARDPICKER_START_Y*row))
		op3.GeoM.Translate(float64(CARDPICKER_START_X+CARDPICKER_DIST_X*col), float64(CARDPICKER_START_Y+CARDPICKER_DIST_Y*row))
		screen.DrawImage(ebitenCard.image, op3)
		cardList = append(cardList, ebitenCard)
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		xCurInt, yCurInt := ebiten.CursorPosition()
		//fmt.Println("DDDDDD", xCur, yCur)
		xCur, yCur := float64(xCurInt), float64(yCurInt)
		for _, ec := range cardList {
			// fmt.Println(ec.x, ec.y)
			if xCur > ec.x && xCur < ec.x+math.Floor(ORI_CARD_WIDTH*HAND_SCALE) &&
				yCur > ec.y && yCur < ec.y+math.Floor(ORI_CARD_HEIGHT*HAND_SCALE) {
				c.selectedCard = ec
				// fmt.Println("Sel", c.selectedCard)
			}
		}
		// check if OK button is clicked
		if xCur > CARDPICKER_START_X && xCur < CARDPICKER_START_X+190 &&
			yCur > 540 && yCur < 540+49 && c.selectedCard != nil {
			fmt.Println("Click OK", len(c.cards))
			for idx, j := range c.cards {
				if j == c.selectedCard.card {
					fmt.Println("Send stuff", idx)
					c.pickedCards <- idx
					c.m.currentSubState = c.m.mainState
					c.selectedCard = nil

					//close(c.pickedCards)
					break

				}
			}
		}

		// check if CANCEL button is clicked
		if c.optional && xCur > CARDPICKER_START_X+200 && xCur < CARDPICKER_START_X+200+190 &&
			yCur > 540 && yCur < 540+49 {
			c.pickedCards <- -1
			c.m.currentSubState = c.m.mainState
			c.selectedCard = nil
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		xCurInt, yCurInt := ebiten.CursorPosition()
		//fmt.Println("DDDDDD", xCur, yCur)
		xCur, yCur := float64(xCurInt), float64(yCurInt)
		for i := len(cardList) - 1; i >= 0; i-- {
			ec := cardList[i]
			// fmt.Println(ec.x, ec.y)
			if xCur > ec.x && xCur < ec.x+math.Floor(ORI_CARD_WIDTH*HAND_SCALE) &&
				yCur > ec.y && yCur < ec.y+math.Floor(ORI_CARD_HEIGHT*HAND_SCALE) {
				c.m.detailViewCard = ec
				break
				// fmt.Println("Sel", c.selectedCard)
			}
		}
		c.m.detailState.prevSubState = c
		c.m.currentSubState = c.m.detailState
	}

	if c.selectedCard != nil {
		op3.GeoM.Reset()
		op3.GeoM.Translate(CARDPICKER_START_X, 540)
		screen.DrawImage(c.m.btn, op3)
		text.Draw(screen, "OK", mplusNormalFont, CARDPICKER_START_X+70, 570, color.White)
		op3.GeoM.Reset()
		// op3.GeoM.Scale(4, 4)
		op3.GeoM.Translate(float64(c.selectedCard.x), float64(c.selectedCard.y))
		screen.DrawImage(c.m.checkMark, op3)
	}
	if c.optional {
		op3.GeoM.Reset()
		op3.GeoM.Translate(CARDPICKER_START_X+200, 540)
		screen.DrawImage(c.m.btn, op3)
		text.Draw(screen, "CANCEL", mplusNormalFont, CARDPICKER_START_X+250, 570, color.White)
	}
	// fmt.Println("===")
}
func (c *cardListState) Update() error {
	return nil
}
func (c *cardListState) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(c.m.bgImage2, op)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(1.3, 1.0)
	op2.GeoM.Translate(120, 0)
	screen.DrawImage(c.m.paperBg, op2)
	op3 := &ebiten.DrawImageOptions{}
	colPerRow := 7
	cardList := []*EbitenCard{}
	for idx, cc := range c.cards {
		ebitenCard := NewEbitenCardFromCard(cc)
		op3.GeoM.Reset()
		op3.GeoM.Scale(HAND_SCALE, HAND_SCALE)
		col := (idx % colPerRow)
		row := (idx / colPerRow)
		// fmt.Println(row, col)
		ebitenCard.x = float64(CARDPICKER_START_X + CARDPICKER_DIST_X*col)
		ebitenCard.y = float64(CARDPICKER_START_Y + CARDPICKER_DIST_Y*row)
		// op3.GeoM.Translate(float64(CARDPICKER_START_X*col), float64(CARDPICKER_START_Y*row))
		op3.GeoM.Translate(float64(CARDPICKER_START_X+CARDPICKER_DIST_X*col), float64(CARDPICKER_START_Y+CARDPICKER_DIST_Y*row))
		screen.DrawImage(ebitenCard.image, op3)
		cardList = append(cardList, ebitenCard)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		xCurInt, yCurInt := ebiten.CursorPosition()
		//fmt.Println("DDDDDD", xCur, yCur)
		xCur, yCur := float64(xCurInt), float64(yCurInt)
		for i := len(cardList) - 1; i >= 0; i-- {
			ec := cardList[i]
			// fmt.Println(ec.x, ec.y)
			if xCur > ec.x && xCur < ec.x+math.Floor(ORI_CARD_WIDTH*HAND_SCALE) &&
				yCur > ec.y && yCur < ec.y+math.Floor(ORI_CARD_HEIGHT*HAND_SCALE) {
				c.m.detailViewCard = ec
				break
				// fmt.Println("Sel", c.selectedCard)
			}
		}
		c.m.detailState.prevSubState = c
		c.m.currentSubState = c.m.detailState
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		xCurInt, yCurInt := ebiten.CursorPosition()
		//fmt.Println("DDDDDD", xCur, yCur)
		xCur, yCur := float64(xCurInt), float64(yCurInt)
		if xCur > CARDPICKER_START_X+200 && xCur < CARDPICKER_START_X+200+190 &&
			yCur > 540 && yCur < 540+49 {
			// c.pickedCards <- -1
			c.m.currentSubState = c.m.mainState
			// c.selectedCard = nil
		}
	}
	op3.GeoM.Reset()
	op3.GeoM.Translate(CARDPICKER_START_X+200, 540)
	screen.DrawImage(c.m.btn, op3)
	text.Draw(screen, "CANCEL", mplusNormalFont, CARDPICKER_START_X+250, 570, color.White)
}

type boolPickState struct {
	m            *MainGameState
	msg          string
	pickedOption chan (bool)
}

func (c *boolPickState) Update() error {
	return nil
}
func (c *boolPickState) BoolPick(message string) bool {
	c.msg = message
	c.m.mutex.Lock()
	c.m.currentSubState = c
	c.m.mutex.Unlock()
	pickedCards := <-c.pickedOption
	// fmt.Println("Dapat hasil", pickedCards)
	return pickedCards
}
func (c *boolPickState) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	screen.DrawImage(c.m.bgImage2, op)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(1.8, 0.3)
	op2.GeoM.Translate(-50, 200)
	screen.DrawImage(c.m.paperBg, op2)
	text.Draw(screen, c.msg, mplusNormalFont, 100, 300, color.Black)

	// ok button
	op.GeoM.Reset()
	op.GeoM.Translate(700, 400)
	screen.DrawImage(c.m.btn, op)

	// cancel button
	op.GeoM.Reset()
	op.GeoM.Translate(900, 400)
	screen.DrawImage(c.m.btn, op)

	text.Draw(screen, "Yes", mplusNormalFont, 780, 430, color.White)
	text.Draw(screen, "No", mplusNormalFont, 980, 430, color.White)

	// detect click
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		curX, curY := ebiten.CursorPosition()
		fmt.Println("cursors bool pick", curX, curY)
		if curY > 400 && curY < 450 {
			if curX > 900 {
				c.pickedOption <- false
			} else if curX > 700 {
				c.pickedOption <- true
			}
			c.m.currentSubState = c.m.mainState
		}
	}

}
