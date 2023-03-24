package main

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
)

func NewEbitenCardFromCard(c cards.Card) *EbitenCard {
	cardImageName := fmt.Sprintf("%s.png", c.GetName())
	//fmt.Println(cardImageName)
	imageProvider := NewImageProvider()
	cardImage := imageProvider.GetImage(cardImageName)
	h := &EbitenCard{image: cardImage, card: c, oriWidth: ORI_CARD_WIDTH, oriHeight: ORI_CARD_HEIGHT}
	return h
}
