package cards

// composite action, an action consisting of smaller sub-action, executed orderly
type CompositeAction struct {
	s       AbstractGamestate
	actions []AbstractActon
}

func (c *CompositeAction) DoAction() {
	for _, v := range c.actions {
		v.DoAction()
	}
}

func NewCompositeAction(s AbstractGamestate, actions ...AbstractActon) AbstractActon {
	comp := &CompositeAction{s: s}
	comp.actions = []AbstractActon{}
	comp.actions = append(comp.actions, actions...)
	return comp
}
