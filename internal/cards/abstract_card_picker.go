package cards

// this interface shows some UI to pick cards
type AbstractCardPicker interface {
	// show message and options, return the index of the picked card
	// list is the list of card where user can pick
	// message is the message
	PickCard(list []Card, message string) int
}
