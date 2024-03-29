package cards

import (
	"fmt"
	"math/rand"
)

type TombForgottenMonarchEntrance struct {
	BaseArea
	state AbstractGamestate
}

func NewTombMonarchEntrance(state AbstractGamestate) TombForgottenMonarchEntrance {
	h := TombForgottenMonarchEntrance{state: state}
	return h
}

func (a *TombForgottenMonarchEntrance) GetName() string {
	return "TombOfForgottenMonachEntrance"
}
func (a *TombForgottenMonarchEntrance) GetDescription() string {
	return "Rewards: 100 Money, Release undead monsters to center deck on explore"
}
func (a *TombForgottenMonarchEntrance) GetCost() Cost {
	cost := NewCost()
	cost.AddResource(RESOURCE_NAME_EXPLORATION, 3)
	return cost
}
func (a *TombForgottenMonarchEntrance) Unbanishable() {

}
func (a *TombForgottenMonarchEntrance) OnExplored() {
	// add reward
	a.state.AddResource(RESOURCE_NAME_MONEY, 100)

	// add progress to the game
	skeletonGuardCount := 3 //(rand.Int() % 6) + 3
	allAddedCard := []Card{}
	for i := 0; i < skeletonGuardCount; i++ {
		k := NewSkeletonGuard(a.state)
		allAddedCard = append(allAddedCard, &k)
	}
	fmt.Println("Len of skeletonguard", skeletonGuardCount)
	// add defeatSkeletonGuardCountEventListener
	// defeat all skeleton guard to add lich mage to center deck
	aa := NewSkeletonGuard(a.state)
	cardFilter := &CardFilter{Key: FILTER_NAME, Op: Eq, Value: aa.GetName()}

	lichMageMonster := NewLichMageMonster(a.state)
	cardsAdded := []Card{&lichMageMonster}
	pushCenterDeckAction := NewPushCenterDeckAction(a.state, cardsAdded, false)
	removeEventListenerAction := NewRemoveEventListenerAction(a.state, EVENT_CARD_DEFEATED, nil)
	compositeAction := NewCompositeAction(a.state, pushCenterDeckAction, removeEventListenerAction)
	countDownAction := NewCountDownAction(2, 1, compositeAction)
	skeletonGuardDefeatedListener := NewCardDefeatedListener(cardFilter, countDownAction)
	removeEventListenerAction.(*RemoveEventListenerAction).listener = append(removeEventListenerAction.(*RemoveEventListenerAction).listener, skeletonGuardDefeatedListener)
	a.state.AttachListener(EVENT_CARD_DEFEATED, skeletonGuardDefeatedListener)
	// add some combat generating card to ensure we can defeat mini boss
	atk := rand.Int() % 3
	switch atk {
	case 0:
		j := NewArcher(a.state)
		allAddedCard = append(allAddedCard, &j)
	case 1:
		j := NewMonsterSlayer(a.state)
		allAddedCard = append(allAddedCard, &j)
	case 2:
		j := NewArcher(a.state)
		allAddedCard = append(allAddedCard, &j)
	}
	a.state.AddCardToCenterDeck(DISCARD_SOURCE_NAN, true, allAddedCard...)

}
