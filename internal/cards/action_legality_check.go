package cards

import "fmt"

type LegalChecker interface {
	Check(interface{}) bool
}
type RuleEnforcer struct {
	rules []LegalChecker
}

func NewRuleEnforcer() *RuleEnforcer {
	return &RuleEnforcer{rules: []LegalChecker{}}
}
func (l *RuleEnforcer) Len() int {
	if l.rules == nil {
		return 0
	}
	return len(l.rules)
}
func (l *RuleEnforcer) Check(data interface{}) bool {
	fmt.Println("Check RuleEnforcer", len(l.rules))
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
	fmt.Println("Attach Rule")
	l.rules = append(l.rules, k)
}
func (l *RuleEnforcer) DetachRule(k LegalChecker) {
	//l.rules = append(l.rules, k)
	fmt.Println("Detach Rule")
	idx := -1
	for i := 0; i < len(l.rules); i++ {
		fmt.Println("TTT", l.rules[i], k, l.rules[i] == k)
		if l.rules[i] == k {
			idx = i
			break
		}
	}
	fmt.Println(idx)
	if idx == 0 {
		l.rules = []LegalChecker{}
	}
	if idx > 0 {
		ll := l.rules[:idx-1]
		newRules := append(ll, l.rules[idx+1:]...)
		l.rules = newRules
	}

}
