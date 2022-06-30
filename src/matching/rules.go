package matching

import (
	"bytes"
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
	return &rr
}

func NewDefaultRuleRegistry() (*RuleRegistry, error) {
	r := NewRuleRegistry().
		Register("startsWith", "^").
		Register("endsWith", "$")

	return r, r.Error
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
	return rr.RegisterWithRegistry(ruleName, baseRule, addedExpression, rr)
}

func (rr *RuleRegistry) RegisterWithBefore(ruleName RuleName, baseRule RuleName, addedExpression Regex) *RuleRegistry {
	return rr.RegisterWithRegistryBefore(ruleName, baseRule, addedExpression, rr)
}

func (rr *RuleRegistry) RegisterWithRegistry(ruleName RuleName, baseRule RuleName, addedExpression Regex, registry *RuleRegistry) *RuleRegistry {
	newRule := registry.rules[baseRule].Regex + addedExpression
	return rr.Register(ruleName, newRule)
}

func (rr *RuleRegistry) RegisterWithRegistryBefore(ruleName RuleName, baseRule RuleName, addedExpression Regex, registry *RuleRegistry) *RuleRegistry {
	newRule := addedExpression + registry.rules[baseRule].Regex
	return rr.Register(ruleName, newRule)
}


type RuleBuilder struct {
	rule RuleName
	buffer bytes.Buffer
	Error error
}

func NewRuleBuilder(rule RuleName) *RuleBuilder {
	return &RuleBuilder{rule: rule}
}

func (rb *RuleBuilder) writeIfNoError(str string) *RuleBuilder {
	if _, err := rb.buffer.WriteString(str); err != nil {
		rb.Error = err
	}
	return rb
}

func (rb *RuleBuilder) With(rule RuleName, str string) *RuleBuilder {
	if rb.Error != nil {
		if rb.writeIfNoError(string(rule)); rb.Error == nil {
			rb.writeIfNoError(str)
		}
	}
	return rb
}

func (rb *RuleBuilder) WithUnregistered(str string) *RuleBuilder {
	if rb.Error != nil {
		rb.writeIfNoError(str)
	}
	return rb
}

func (rb *RuleBuilder) Build() string {
	return rb.buffer.String()
}
