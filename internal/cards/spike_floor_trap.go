package cards

type SpikeFloor struct {
	BaseTrap
	state      AbstractGamestate
	isDisarmed bool
}

func NewSpikeFloor(state AbstractGamestate) SpikeFloor {
	return SpikeFloor{state: state}
}
func (b *SpikeFloor) GetName() string {
	return "SpikeFloor"
}
func (b *SpikeFloor) GetDescription() string {
	return "when enter center row: take 4 damage"
}
func (b *SpikeFloor) Dispose(source string) {
	b.state.BanishCard(b, DISCARD_SOURCE_CENTER)
}

func (b *SpikeFloor) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 2)
	return cost
}

func (b *SpikeFloor) Trap() {
	if !b.isDisarmed {
		b.state.TakeDamage(4)
	}
}
func (b *SpikeFloor) IsDisarmed() bool {
	return b.isDisarmed
}
func (b *SpikeFloor) Disarm() {
	b.isDisarmed = true
}
func (b *SpikeFloor) OnDisarm() {

}
