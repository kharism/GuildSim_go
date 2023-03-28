package cards

type BaseItem struct{}

func (h *BaseItem) GetName() string {
	return ""
}
func (h *BaseItem) GetDescription() string {
	return ""
}
func (h *BaseItem) GetCost() Cost {
	cost := NewCost()
	return cost
}

// return Hero,Area,Monster,Event etc
func (h *BaseItem) GetCardType() CardType {
	return Item
}
func (h *BaseItem) OnAddedToHand() {}

// when played from hand, do this
func (h *BaseItem) OnPlay() {}

func (h *BaseItem) OnAcquire() {}

// when explored, do this
func (h *BaseItem) OnExplored() {}

// when slain, do this
func (h *BaseItem) OnSlain() {}

// when discarded to cooldown pile, do this
func (h *BaseItem) OnDiscarded() {}

func (h *BaseItem) Dispose(source string) {

}

// specific type of item. implement this if you want to implement 'usable' item
// that is not provide passive buff/debuff
// the item can be discarded to cooldown or banished or just stay on the item slot
type Consumable interface {
	OnConsume()
}
