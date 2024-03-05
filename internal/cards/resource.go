package cards

import (
	"errors"
	"fmt"
	"strings"
)

var RESOURCE_NOT_FOUND_ERROR = errors.New("Resource not found")

var RESOURCE_NAME_EXPLORATION = "Exploration"
var RESOURCE_NAME_COMBAT = "Combat"
var RESOURCE_NAME_MONEY = "Money"
var RESOURCE_NAME_REPUTATION = "Reputation"
var RESOURCE_NAME_BLOCK = "Block"

// resource pool for the game
type Resource struct {
	Detail map[string]int
}

func NewResource() Resource {
	res := Resource{}
	res.Detail = map[string]int{}
	return res
}

// it's add resource
// param key is resource name
// param i is the number, only add positive integer
func (r *Resource) AddResource(key string, i int) {
	if i < 0 {
		return
	}
	r.Detail[key] += i
}

// remove resource
// param key is resource name
// param i is number of resource removed, use positive integer only
// if no resource is present return RESOURCE_NOT_FOUND_ERROR

func (r *Resource) RemoveResource(key string, i int) error {
	if _, ok := r.Detail[key]; !ok {
		return RESOURCE_NOT_FOUND_ERROR
	}
	if r.Detail[key] >= i {
		r.Detail[key] -= i
	} else {
		r.Detail[key] = 0
	}

	return nil
}

type Cost struct {
	*Resource
}

func NewCost() Cost {
	cost := Cost{}
	res := NewResource()
	cost.Resource = &res
	cost.Detail = map[string]int{}
	return cost
}

func (c Cost) String() string {
	output := []string{}

	for key, val := range c.Detail {
		output = append(output, fmt.Sprintf("%s:%d", key, val))
	}
	return strings.Join(output, "|")
}

// compare the cost to resource
func (c *Cost) IsEnough(res Resource) bool {
	if len(c.Detail) == 0 {
		return true
	}
	for key, val := range c.Detail {
		if _, ok := res.Detail[key]; !ok {
			return false
		}
		if val > res.Detail[key] {
			return false
		}
	}
	return true
}
