package cards

var Eq = "eq"
var Neq = "neq"
var Gt = "gt"
var Lt = "lt"

var FILTER_NAME = "NAME"

// var FILTER_NAME = "NAME"
type CardFilter struct {
	Key   string
	Op    string
	Value interface{}
}

func Match(c Card, filter CardFilter) bool {
	if filter.Key == FILTER_NAME {
		val_str := filter.Value.(string)
		if filter.Op == Eq && val_str == c.GetName() {
			return true
		} else if filter.Op == Neq && val_str != c.GetName() {
			return true
		}
	}
	return false
}
func Contains(pile []Card, filter CardFilter) bool {
	for _, c := range pile {
		if Match(c, filter) {
			return true
		}
	}
	return false
}
