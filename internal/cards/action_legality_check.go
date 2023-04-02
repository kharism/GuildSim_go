package cards

type LegalChecker interface {
	Check(interface{}) bool
}
type RuleEnforcer struct {
	rules []LegalChecker
}

func NewRuleEnforcer() *RuleEnforcer {
	return &RuleEnforcer{rules: []LegalChecker{}}
}
func (l *RuleEnforcer) Check(data interface{}) bool {
	if len(l.rules) == 0 {
		return true
	}
	output := true
	for _, rule := range l.rules {
		output = output && rule.Check(data)
	}
	return output
}
func (l *RuleEnforcer) AttachRule(k LegalChecker) {
	l.rules = append(l.rules, k)
}
func (l *RuleEnforcer) DetachRule(k LegalChecker) {
	//l.rules = append(l.rules, k)
	idx := -1
	for i := 0; i < len(l.rules); i++ {
		if l.rules[i] == k {
			idx = i
			break
		}
	}
	if idx == 0 {
		l.rules = []LegalChecker{}
	}
	if idx > 0 {
		ll := l.rules[:idx-1]
		newRules := append(ll, l.rules[idx+1:]...)
		l.rules = newRules
	}

}
