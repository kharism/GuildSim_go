package cards

// a specific type of cards that have cards on top of it
type Overlay interface {
	Card
	HasOverlayCard() bool
	AttachOverlayCard(cards Card)
	GetOverlay() []Card
	// detach the top card of this stack
	Detach()
}
