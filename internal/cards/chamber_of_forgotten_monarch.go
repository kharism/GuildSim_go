package cards

type ChamberOfForgottenMonarch struct {
	BaseArea
	state AbstractGamestate
}

func NewChamberOfForgottenMonarch(state AbstractGamestate) ChamberOfForgottenMonarch {
	a := ChamberOfForgottenMonarch{state: state}
	return a
}
func (ed *ChamberOfForgottenMonarch) GetName() string {
	return "ChamberOfForgottenMonarch"
}
func (ed *ChamberOfForgottenMonarch) GetDescription() string {
	return "if you have Crystal Key and Bone key, explore this card for free. On explored: stack The forgotten " +
		"monarch to center deck"
}
func (a *ChamberOfForgottenMonarch) GetCost() Cost {
	cost := NewCost()
	items := a.state.ListItems()
	crystalKey := CrystalKey{}
	boneKey := BoneKey{}
	boneKeyFound := false
	crystalKeyFound := false
	for _, i := range items {
		if i.GetName() == crystalKey.GetName() {
			crystalKeyFound = true
		}
		if i.GetName() == boneKey.GetName() {
			boneKeyFound = true
		}
	}
	if !boneKeyFound || !crystalKeyFound {
		cost.AddResource(RESOURCE_NAME_EXPLORATION, 99)
	}

	return cost
}
func (a *ChamberOfForgottenMonarch) Unbanishable() {}
func (a *ChamberOfForgottenMonarch) OnExplored() {
	forgottenMonarch := NewForgottenMonarch(a.state)
	a.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, false, &forgottenMonarch)
}
