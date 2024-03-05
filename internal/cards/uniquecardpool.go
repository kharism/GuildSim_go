package cards

import "math/rand"

type UniqueCardPool struct {
	MonsterPool []Card
}

var pool *UniqueCardPool

func GetPool(state AbstractGamestate) *UniqueCardPool {
	if pool == nil {
		curseEater := NewCurseEaterMonster(state)
		potionVendor := NewPotionVendor(state)
		pool = &UniqueCardPool{
			MonsterPool: []Card{
				&curseEater,
				NewRagingVampire(state),
				&potionVendor,
			},
		}
		rand.Shuffle(len(pool.MonsterPool), func(i, j int) {
			pool.MonsterPool[i], pool.MonsterPool[j] = pool.MonsterPool[j], pool.MonsterPool[i]
		})
	}
	return pool
}
func ResetPool() {
	pool = nil
}
func (p *UniqueCardPool) Fetch() Card {
	o := p.MonsterPool[0]
	p.MonsterPool = p.MonsterPool[1:]
	return o
}
