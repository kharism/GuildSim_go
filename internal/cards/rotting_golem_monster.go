package cards

type RottingGolem struct {
	BaseMonster
	state AbstractGamestate
}
type BoneKey struct {
	BaseItem
}

func (h *BoneKey) GetName() string {
	return "BoneKey"
}
func (h *BoneKey) GetDescription() string {
	return "Open path to forgotten monarch"
}

func NewRottingGolem(a AbstractGamestate) RottingGolem {
	return RottingGolem{state: a}
}
func (t *RottingGolem) GetName() string {
	return "RottingGolem"
}
func (m *RottingGolem) GetDescription() string {
	return "stack carcass slime curse,  on defeat, gains 1 key item"
}
func (m *RottingGolem) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 7)
	return cost
}
func (d *RottingGolem) OnSlain() {
	boneKey := BoneKey{}
	d.state.AddItem(&boneKey)
}
func (d *RottingGolem) Unshuffleable() {}
func (d *RottingGolem) Unbanishable()  {}
func (d *RottingGolem) OnPunish() {
	//d.state.TakeDamage(8)
	curse := NewCarcassSlimeCurse(d.state)
	d.state.StackCards(DISCARD_SOURCE_NAN, &curse)
}
