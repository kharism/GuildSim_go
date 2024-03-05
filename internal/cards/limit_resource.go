package cards

type LimitResource struct {
	state           AbstractGamestate
	resourceReducer *ResourceReducerAction
}

// this action check if a player generate a resource that is equal or grater thatn threshold
// it then reduce the resource by ReduceAmount
type ResourceReducerAction struct {
	state        AbstractGamestate
	ResourceName string
	Threshold    int
	ReduceAmount int
}

func (r *ResourceReducerAction) DoAction(data map[string]interface{}) {
	resourceName := data[EVENT_ATTR_ADD_RESOURCE_NAME].(string)
	resourceAmount := data[EVENT_ATTR_ADD_RESOURCE_AMOUNT].(int)
	if resourceName == r.ResourceName && resourceAmount > 0 && resourceAmount >= r.Threshold {
		r.state.AddResource(resourceName, -1)
	}

}

// Create a resourceReducer struct, it has an action and
func NewLimitResource(state AbstractGamestate, resourceName string, threshold, reducer int) *LimitResource {
	reducerAction := &ResourceReducerAction{ResourceName: resourceName, Threshold: threshold, ReduceAmount: reducer}
	return &LimitResource{state: state, resourceReducer: reducerAction}
}

func (r *LimitResource) AttachListener() {
	r.state.AttachListener(EVENT_ADD_RESOURCE, r.resourceReducer)
}
func (r *LimitResource) DetachListener() {
	r.state.RemoveListener(EVENT_ADD_RESOURCE, r.resourceReducer)
}
