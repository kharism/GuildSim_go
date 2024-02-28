package cards

type RagingVampire struct {
	BaseMonster
	state         AbstractGamestate
	bleedMechanic *BleedingMechanic
	isPlayed      bool
}
type BleedingMechanic struct {
	source Card
	state  AbstractGamestate
}

func (b *BleedingMechanic) DoAction(data map[string]interface{}) {
	playedCard := data[EVENT_ATTR_CARD_PLAYED].(Card)
	// fmt.Println(&playedCard)
	// fmt.Println(&b.source)
	if playedCard != b.source {
		b.state.TakeDamage(1)
	}

}
func NewRagingVampire(state AbstractGamestate) Card {
	res := RagingVampire{state: state}
	bleedMechanic := &BleedingMechanic{state: state, source: &res}
	res.bleedMechanic = bleedMechanic
	return &res
}

func (b *RagingVampire) GetName() string {
	return "RagingVampire"
}
func (b *RagingVampire) GetDescription() string {
	return "recruitable. On punish: take 1 damage. On recruit: loose 10 reputation. On play: if it is the fist card you played, gain 8 combat but take" +
		" 1 damage each time you play another card"
}
func (b *RagingVampire) Dispose(source string) {
	b.state.RemoveListener(EVENT_CARD_PLAYED, b.bleedMechanic)
	if b.isPlayed {
		b.state.DiscardCard(b, DISCARD_SOURCE_PLAYED)
	} else {
		b.state.DiscardCard(b, DISCARD_SOURCE_HAND)
	}
	b.isPlayed = false
}
func (b *RagingVampire) GetKeywords() []string {
	return []string{"Bleed: take 1 damage each time \nyou play another card this turn"}
}
func (b *RagingVampire) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_COMBAT, 2)
	return cost
}
func (r *RagingVampire) OnPunish() {
	r.state.TakeDamage(1)
}
func (r *RagingVampire) OnSlain() {

}
func (r *RagingVampire) OnRecruit() {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_REPUTATION, 10)
	if r.state.GetCurrentResource().Detail[RESOURCE_NAME_REPUTATION] >= 10 {
		r.state.PayResource(cost)
	} else {
		r.state.GetCurrentResource().Detail[RESOURCE_NAME_REPUTATION] = 0
	}

}
func (r *RagingVampire) OnPlay() {
	if len(r.state.GetPlayedCards()) == 0 {
		r.state.AddResource(RESOURCE_NAME_COMBAT, 10)
	}
	r.isPlayed = true
	r.state.AttachListener(EVENT_CARD_PLAYED, r.bleedMechanic)
}
