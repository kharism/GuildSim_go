package cards

type SkeletonGuard struct {
	BaseMonster
	state AbstractGamestate
}

func NewSkeletonGuard(state AbstractGamestate) SkeletonGuard {
	k := SkeletonGuard{state: state}
	return k
}
func (m *SkeletonGuard) GetName() string {
	return "Skeleton Guard"
}
func (m *SkeletonGuard) GetDescription() string {
	return "Deals 3 dmg, on slain gains 1 exploration and 1 reputation"
}
func (m *SkeletonGuard) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 3)
	return cost
}
func (m *SkeletonGuard) OnPunish() {
	m.state.TakeDamage(3)
}
func (m *SkeletonGuard) OnSlain() {
	m.state.AddResource(RESOURCE_NAME_EXPLORATION, 1)
	m.state.AddResource(RESOURCE_NAME_REPUTATION, 1)
}
