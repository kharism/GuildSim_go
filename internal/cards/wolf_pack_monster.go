package cards

type WolfPack struct {
	BaseMonster
	state       AbstractGamestate
	turnCounter int
	isDisarmed  bool
}

func (b *WolfPack) GetName() string {
	return "WolfPack"
}
func (b *WolfPack) GetDescription() string {
	return "trap: stack 2 copies of DireWolf to central deck. On punish: 2 dmg on end of turn. rewards: 2 reputation and choose either 3 exploration point or draw a card"
}
func (b *WolfPack) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}

func (b *WolfPack) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 4)
	return cost
}

func (b *WolfPack) OnSlain() {
	b.state.AddResource(RESOURCE_NAME_REPUTATION, 2)
	if b.state.GetBoolPicker().BoolPick("Draw a card?") {
		b.state.Draw()
	} else {
		b.state.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	}
}

func NewWolfPack(state AbstractGamestate) WolfPack {
	k := WolfPack{state: state, isDisarmed: false}
	return k
}

func (b *WolfPack) Trap() {
	if !b.isDisarmed {
		for i := 0; i < 2; i++ {
			direwolf := NewDirewolf(b.state)
			b.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, false, &direwolf)
		}
	}
}
func (b *WolfPack) OnDisarm() {
	b.OnSlain()
}
func (b *WolfPack) IsDisarmed() bool {
	return b.isDisarmed
}
func (b *WolfPack) Disarm() {
	b.isDisarmed = true
}
