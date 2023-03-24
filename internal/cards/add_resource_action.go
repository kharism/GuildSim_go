package cards

type AddResourceAction struct {
	s      AbstractGamestate
	n      string
	amount int
}

// DoAction implements observer.Listener
func (a *AddResourceAction) DoAction() {
	// fmt.Println("Add Resource", a.n, a.amount)
	a.s.AddResource(a.n, a.amount)
}

func NewAddResourceAction(s AbstractGamestate, resourceName string, amount int) *AddResourceAction {
	dd := &AddResourceAction{s: s}
	dd.n = resourceName
	dd.amount = amount
	return dd
}
