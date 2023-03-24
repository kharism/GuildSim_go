package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ImageProvide provide image, duh
// it also cache the card image so we don't load the image to memory each time we load a card to screen
type ImageProvider struct {
	cardCache map[string]*ebiten.Image
}

var imgProvider *ImageProvider

func NewImageProvider() *ImageProvider {
	if imgProvider != nil {
		return imgProvider
	}
	imgProvider = &ImageProvider{cardCache: make(map[string]*ebiten.Image)}
	return imgProvider
}

func (ip *ImageProvider) GetImage(filename string) *ebiten.Image {
	if _, ok := ip.cardCache[filename]; ok {
		return ip.cardCache[filename]
	}
	cardImg, _, err := ebitenutil.NewImageFromFile("img/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	ip.cardCache[filename] = cardImg
	return cardImg
}
