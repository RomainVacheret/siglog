package matching

import (
	"regexp"
	log "github.com/sirupsen/logrus"
)

type RuleName string

type Regex string

type Rule struct {
	RuleName RuleName
	Regex Regex
	RegexExpr *regexp.Regexp
}

type RuleRegistry struct {
	rules map[RuleName]Rule
	Error error
}

func NewRuleRegistry() *RuleRegistry {
	rr := RuleRegistry{rules: make(map[RuleName]Rule)}
	rr.initialize()
	return &rr
}

func (rr *RuleRegistry) initialize() error {
	rr.Register("startsWith", "^")
	rr.Register("endsWith", "$")
	return rr.Error
}

func (rr *RuleRegistry) Register(ruleName RuleName, rule Regex) *RuleRegistry {
	if regex, err := regexp.Compile(string(rule)); err == nil {
		rr.rules[ruleName] = Rule{
			RuleName: ruleName,
			Regex: rule,
			RegexExpr: regex,
		}
		
		return rr
	} else {
		log.WithFields(log.Fields{
			"ruleName": ruleName,
			"regex": regex,
			"err": err,
		}).Error("An error occured while registering a new rule")
		rr.Error = err

		return rr
	}
}

func (rr *RuleRegistry) RegisterWith(ruleName RuleName, baseRule RuleName, addedExpression Regex) *RuleRegistry {
	newRule := rr.rules[baseRule].Regex + addedExpression
	return rr.Register(ruleName, newRule)
}

func (rr *RuleRegistry) RegisterWithBefore(ruleName RuleName, baseRule RuleName, addedExpression Regex) *RuleRegistry {
	newRule := addedExpression + rr.rules[baseRule].Regex
	return rr.Register(ruleName, newRule)
}