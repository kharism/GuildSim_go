package cards

type SkeletonMinion struct {
	BaseMonster
	state AbstractGamestate
}

func NewSkeletonMinion(state AbstractGamestate) SkeletonMinion {
	k := SkeletonMinion{state: state}
	return k
}
func (m *SkeletonMinion) GetName() string {
	return "SkeletonMinion"
}
func (m *SkeletonMinion) GetDescription() string {
	return "deals 3 damage"
}
func (m *SkeletonMinion) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 4)
	return cost
}
func (m *SkeletonMinion) OnPunish() {
	m.state.TakeDamage(3)
}
func (m *SkeletonMinion) OnSlain() {

}
