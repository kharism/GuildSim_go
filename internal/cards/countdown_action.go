package cards

// an action that do countdown, once counter is less than equal 0, will execute main action
type CountDownAction struct {
	counter    int
	decrement  int
	mainAction AbstractActon
}

func NewCountDownAction(counter, decrement int, action AbstractActon) *CountDownAction {
	return &CountDownAction{counter: counter, decrement: decrement, mainAction: action}
}

func (c *CountDownAction) DoAction() {
	c.counter -= c.decrement
	if c.counter <= 0 {
		c.mainAction.DoAction()
	}
}
