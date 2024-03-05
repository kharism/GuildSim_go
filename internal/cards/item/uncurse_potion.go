package item

import "github/kharism/GuildSim_go/internal/cards"

type UncursePotion struct {
	cards.BaseItem
	state cards.AbstractGamestate
}

func NewUncursePotion(state cards.AbstractGamestate) UncursePotion {
	return UncursePotion{state: state}
}

func (h *UncursePotion) Dispose(source string) {
	h.state.BanishCard(h, source)
}
func (h *UncursePotion) GetName() string {
	return "UncursePotion"
}
func (h *UncursePotion) GetDescription() string {
	return "Banish all curse on hand, then draw that many card"
}
func (h *UncursePotion) GetCost() cards.Cost {
	cost := cards.NewCost()
	return cost
}

func (h *UncursePotion) OnConsume() {
	centerCard := h.state.GetCardInHand()
	curseOnHand := []cards.Card{}
	for _, c := range centerCard {
		if c.GetCardType() == cards.Curse {
			curseOnHand = append(curseOnHand, c)
		}
	}
	length := len(curseOnHand)
	// for _, t := range curseOnHand {
	// 	h.state.BanishCard(t, cards.DISCARD_SOURCE_HAND)
	// }
	for i := len(curseOnHand) - 1; i >= 0; i-- {
		h.state.RemoveCardFromHand(curseOnHand[i])
		h.state.BanishCard(curseOnHand[i], cards.DISCARD_SOURCE_HAND)

	}
	for i := 0; i < length; i++ {
		h.state.Draw()
	}
}
