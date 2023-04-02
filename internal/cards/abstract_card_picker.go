package cards

// this interface shows some UI to pick cards
type AbstractCardPicker interface {
	// show message and options, use this for mandatory select a card
	// return the index of the picked card
	// list is the list of card where user can pick
	// message is the message
	PickCard(list []Card, message string) int

	// the same with PickCard, but this time optional
	// return negative number if no cards are choosen
	PickCardOptional(list []Card, message string) int
}

type AbstractBoolPicker interface {
	// as yes/no question
	BoolPick(message string) bool
}

type AbstractDetailViewer interface {
	ShowDetail(Card)
}
