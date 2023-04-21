package cards

import "fmt"

type ShockCurse struct {
	BaseCurse
	state   AbstractGamestate
	limiter *LimitDraw
}

func NewShockCurse(state AbstractGamestate) ShockCurse {
	ld := NewLimitDraw(state, 5)
	return ShockCurse{state: state, limiter: ld}
}

func (d *ShockCurse) GetName() string {
	return "ShockCurse"
}
func (d *ShockCurse) GetDescription() string {
	return "If drawn: take 3 damage and continuous eff Limitdraw to 5"
}

// when added to hand do this
func (d *ShockCurse) OnAddedToHand() {
	d.state.TakeDamage(3)
	d.limiter.AttachLimitDraw(d.state)
}
func (d *ShockCurse) OnBanished() {
	d.limiter.DetachLimitDraw(d.state)
}
func (d *ShockCurse) OnDiscarded() {
	d.limiter.DetachLimitDraw(d.state)
}

func (d *ShockCurse) OnReturnToDeck() {
	fmt.Println("Detach limiter")
	d.limiter.DetachLimitDraw(d.state)
}

func (d *ShockCurse) Dispose(source string) {
	d.state.DiscardCard(d, source)
}
