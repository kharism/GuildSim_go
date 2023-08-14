package cards

var Eq = "eq"
var Neq = "neq"
var Gt = "gt"
var Lt = "lt"
var In = "in"

var FILTER_NAME = "NAME"
var FILTER_TYPE = "TYPE"

// var FILTER_NAME = "NAME"
type CardFilter struct {
	Key   string
	Op    string
	Value interface{}
}

func Match(c Card, filter *CardFilter) bool {
	if filter.Key == FILTER_NAME {
		val := filter.Value
		if filter.Op == Eq && val.(string) == c.GetName() {
			return true
		} else if filter.Op == Neq && val.(string) != c.GetName() {
			return true
		} else if filter.Op == In {
			names := val.([]string)
			for _, n := range names {
				if n == c.GetName() {
					return true
				}
			}
			return false
		}
	} else if filter.Key == FILTER_TYPE {
		val := filter.Value
		if filter.Op == Eq && c.GetCardType() == val.(CardType) {
			return true
		} else if filter.Op == Neq && val.(CardType) != c.GetCardType() {
			return true
		} else if filter.Op == In {
			cardTypes := val.([]CardType)
			for _, t := range cardTypes {
				if t == c.GetCardType() {
					return true
				}
			}
		}
		return false
	}
	return false
}
func Contains(pile []Card, filter *CardFilter) bool {
	for _, c := range pile {
		if Match(c, filter) {
			return true
		}
	}
	return false
}
