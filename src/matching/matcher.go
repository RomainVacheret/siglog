package matching


type Matcher interface {
	Match(text []string) ([]string, error)
	CountMatches(text []string) (int, error)
}

type matcherUtils struct {
	registry RuleRegistry
}


func NewMatcherUtils(registry RuleRegistry) *matcherUtils {
	return &matcherUtils{registry: registry}
}

func (mu *matcherUtils) MatchString(rule RuleName, str string) bool {
	return mu.registry.rules[rule].RegexExpr.MatchString(str)
}

func (mu *matcherUtils) MatchLines(rule RuleName, lines []string) []int {
	// TODO Size could be reduced
	matches_idx := make([]int, len(lines))
	for idx, line := range lines {
		if(mu.MatchString(rule, line)) {
			matches_idx = append(matches_idx, idx)
		}
	}
	return matches_idx
}

// Return the first rule that matched
func (mu *matcherUtils) MatchAny(str string) (RuleName, bool) {
	for name, _ := range mu.registry.rules {
		if hasMatched := mu.MatchString(name, str); hasMatched == true {
			return name, true
		}
	}

	return "", false
}


type UnorderedMatcher struct {
	registry RuleRegistry
	*matcherUtils
}

func NewUnorderedMatcher (registry RuleRegistry) *UnorderedMatcher {
	return &UnorderedMatcher{
		registry,
		NewMatcherUtils(registry),
	}
}

func (um *UnorderedMatcher) Match(lines []string) ([]string, []RuleName, error) {
	matchedLines := make([]string, len(lines))
	matchingRules := make([]RuleName, len(lines))
	totalLength := 0

	for _, line := range lines {
		if name, hasMatched := um.MatchAny(line); hasMatched {
			matchedLines[totalLength] = line
			matchingRules[totalLength] = name
			totalLength++
		}
	}
	return matchedLines[:totalLength], matchingRules[:totalLength] ,nil
}

func (um *UnorderedMatcher) CountMatches(lines []string) (int, error) {
	if l, _, err := um.Match(lines); err == nil {
		return len(l), err
	} else {
		return 0, err
	}
}




type OrderedMatcher struct {
	registry RuleRegistry
	rules []RuleName
}

func NewOrderedMatcher(registry RuleRegistry, rules []RuleName) *OrderedMatcher {
	return &OrderedMatcher{registry: registry, rules: rules}
}


