package observer

// use observer pattern here
type Observer interface {
	Attach(l Listener)
	Detach(l Listener)
	Notify(data map[string]interface{})
}

type Listener interface {
	DoAction(data map[string]interface{})
}
