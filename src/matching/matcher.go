package matching

type Matcher struct {
	registry RuleRegistry
}

func NewMatcher(registry RuleRegistry) *Matcher {
	return &Matcher{registry: registry}
}

func (m *Matcher) MatchString(rule RuleName, str string) bool {
	return m.registry.rules[rule].RegexExpr.MatchString(str)
}

func (m *Matcher) MatchLines(rule RuleName, lines []string) []int {
	// TODO Size could be reduced
	matches_idx := make([]int, len(lines))
	for idx, line := range lines {
		if(m.MatchString(rule, line)) {
			matches_idx = append(matches_idx, idx)
		}
	}
	return matches_idx
}