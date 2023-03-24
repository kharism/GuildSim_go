package cards

type EasyDungeonArea struct {
	BaseArea
	state AbstractGamestate
}

func NewEasyDungeonArea(state AbstractGamestate) EasyDungeonArea {
	a := EasyDungeonArea{state: state}
	return a
}
func (ed *EasyDungeonArea) GetName() string {
	return "EasyDungeonArea"
}
func (ed *EasyDungeonArea) GetDescription() string {
	return "Reward: 100Money and 1 Reputation"
}
func (ed *EasyDungeonArea) GetCost() Cost {
	c := NewCost()
	c.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return c
}

// when played on hand, to this
func (ed *EasyDungeonArea) OnPlay() {}
func (ed *EasyDungeonArea) OnExplored() {
	ed.state.AddResource(RESOURCE_NAME_MONEY, 100)
	ed.state.AddResource(RESOURCE_NAME_REPUTATION, 1)
}

// when slain, do this
func (ed *EasyDungeonArea) OnSlain() {}

func (ed *EasyDungeonArea) Dispose(source string) {
	ed.state.BanishCard(ed, source)
}

// when discarded to cooldown pile, do this
func (ed *EasyDungeonArea) OnDiscarded() {}
