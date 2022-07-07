package matching

import (
	"errors"
	"fmt"
)

type Matcher interface {
	Match(text []string) ([]string, error)
	CountMatches(text []string) (int, error)
}

type matcherUtils struct {
	registry RuleRegistry
}

type NRuleCategory string 
type nRule struct {
	category NRuleCategory
	n uint
}

const (
	Const NRuleCategory = "x"
	Any = "*"
	Some = '+'
	Maybe = "?"
)


func NewMatcherUtils(registry RuleRegistry) *matcherUtils {
	return &matcherUtils{registry: registry}
}

func (mu *matcherUtils) MatchString(rule RuleName, str string) bool {
	return mu.registry.rules[rule].RegexExpr.MatchString(str)
}

func (mu *matcherUtils) MatchLines(rule RuleName, lines []string) []int {
	// TODO Size could be reduced
	matches_idx := make([]int, len(lines))
	totalMatches := 0
	for idx, line := range lines {
		if(mu.MatchString(rule, line)) {
			matches_idx[totalMatches] = idx
			totalMatches++
		}
	}
	return matches_idx[:totalMatches]
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

// Add nRules
// Add option to match any/all rules for a given line
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
	return matchedLines[:totalLength], matchingRules[:totalLength], nil
}

func (um *UnorderedMatcher) CountMatches(lines []string) (int, error) {
	if l, _, err := um.Match(lines); err == nil {
		return len(l), err
	} else {
		return 0, err
	}
}

func min(a int, b int) int {
	if a < b {
		return a 
	} else {
		return b 
	}
}

type OrderedMatcher struct {
	registry RuleRegistry
	rules []RuleName
	nRules []nRule
	nRulesCount []uint
	*matcherUtils
}

func NewOrderedMatcher(registry RuleRegistry, rules []RuleName) *OrderedMatcher {
	nRules := make([]nRule,len(rules))
	nRulesCount := make([]uint,len(rules))
	for i := 0; i < len(rules); i++ {
		nRules[i] = nRule{Any, 1}
	}
	return &OrderedMatcher{
		registry: registry,
		rules: rules,
		nRules: nRules,
		nRulesCount: nRulesCount,
		matcherUtils: NewMatcherUtils(registry)}
}

func (om *OrderedMatcher) Match(lines []string) ([]string, []RuleName, error) {
	matchedLines := make([]string, len(lines))
	matchingRules := make([]RuleName, len(lines))
	totalLength := 0
	currentRuleIdx := 0

	for _, line := range lines {
		if hasMatched := om.MatchString(om.rules[currentRuleIdx], line); hasMatched {
			fmt.Println(om.rules[currentRuleIdx], line, om.nRules[currentRuleIdx], om.nRulesCount[currentRuleIdx])
			n := om.nRules[currentRuleIdx].n
			
			if n != 0 && om.nRulesCount[currentRuleIdx] ==  n{
				return nil, nil, errors.New("Maximum reached for nRule " + string(om.nRules[currentRuleIdx].category))
			}

			matchedLines[totalLength] = line
			matchingRules[totalLength] = om.rules[currentRuleIdx]
			totalLength++
			om.nRulesCount[currentRuleIdx]++

			if om.nRulesCount[currentRuleIdx] == n {
				currentRuleIdx = min(currentRuleIdx + 1, len(om.rules) - 1)
			}
		} else if currentRuleIdx + 1 < len(om.rules) {
			if hasNextMatched, nNext := om.MatchString(om.rules[currentRuleIdx + 1], line), om.nRules[currentRuleIdx].n; hasNextMatched && nNext != 0 && om.nRulesCount[currentRuleIdx + 1] + 1 < nNext {
			fmt.Println(om.rules[currentRuleIdx + 1], line, "2")
				matchedLines[totalLength] = line
				matchingRules[totalLength] = om.rules[currentRuleIdx + 1]
				om.nRulesCount[currentRuleIdx + 1]++
				totalLength++
				currentRuleIdx = min(currentRuleIdx + 1, len(om.rules) - 1)
			}
		}

	}
	return matchedLines[:totalLength], matchingRules[:totalLength], nil
}


func (om *OrderedMatcher) CountMatches(lines []string) (int, error) {
	if l, _, err := om.Match(lines); err == nil {
		return len(l), err
	} else {
		return 0, err
	}
}
