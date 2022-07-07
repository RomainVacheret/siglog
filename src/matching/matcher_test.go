package matching

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var lines = []string {
	"Welcome traveler",
	"Do you like it here?",
	"The view is wonderful for sure",
	"Warning: the nights are cold",
	"Sometimes, the temperature can be negative",
}


type MatcherUtilsTestSuite struct {
	suite.Suite
	MatcherUtils *matcherUtils
}

func (s *MatcherUtilsTestSuite) SetupTest() {
	defaultRegistry, _ := NewDefaultRuleRegistry()
	registry := NewRuleRegistry()
	s.MatcherUtils = NewMatcherUtils(*registry)

	registry.RegisterWithRegistry("startsWithW", "startsWith", "W", defaultRegistry)
	registry.RegisterWithRegistryBefore("endsWithe", "endsWith", "e", defaultRegistry)
}

func (s *MatcherUtilsTestSuite) TestMatchString() {
	for _, test := range []struct {
		rule RuleName 
		str string
		result bool
	} {
		{rule: "startsWithW", str: "Welcome", result: true},
		{rule: "endsWithe", str: "Welcome", result: true},
		{rule: "startsWithW", str: "Foo", result: false},
		{rule: "endsWithe", str: "Foo", result: false},
		{rule: "startsWithW", str: "here", result: false},
		{rule: "endsWithe", str: "here", result: true},
		{rule: "startsWithW", str: "welcome", result: false},
		{rule: "endsWithe", str: "WelcomE", result: false},
		{rule: "startsWithW", str: "W", result: true},
		{rule: "endsWithe", str: "e", result: true},
	} {
		str := "Equals(\"" + string(test.rule) + ", " + test.str + "\") should be " + strconv.FormatBool(test.result)
		assert.Equal(s.T(), test.result, s.MatcherUtils.MatchString(test.rule, test.str),  str)
	}
}

func (s *MatcherUtilsTestSuite) TestMatchLines() {
	for _, test := range []struct {
		rule RuleName
		result []int
	} {
		{rule: "startsWithW", result: []int{0, 3}},
		{rule: "endsWithe", result: []int{2, 4}},
	} {
		resultAsString := make([]string, len(test.result))
		for idx, val := range test.result {
			resultAsString[idx] = strconv.Itoa(val)
		}
		str := "Equals(\"" + string(test.rule) + ", " + strings.Join(lines, "\n") + "\") should be " + strings.Join(resultAsString, ", ")
		assert.Equal(s.T(), test.result, s.MatcherUtils.MatchLines(test.rule, lines), str)
	}
}

func (s *MatcherUtilsTestSuite) TestMatchAny() {
	for _, test := range[] struct {
		str string
		resultRule RuleName
		resultBool bool
	} {
		{str: "welcomE", resultRule: "", resultBool: false},
		{str: "Welcome", resultRule: "startsWithW", resultBool: true},
		{str: "", resultRule: "", resultBool: false},
		{str: "W", resultRule: "startsWithW", resultBool: true},
		{str: "e", resultRule: "endsWithe", resultBool: true},
	} {
		rules := s.MatcherUtils.registry.rules
		rulesAsString := make([]string, len(rules))
		currentRuleIdx := 0

		for key := range rules {
			rulesAsString[currentRuleIdx] = string(key)
			currentRuleIdx++
		}
		strRule := "Equals(\"" + test.str + strings.Join(rulesAsString, ", ") + "\") should be " + string(test.resultRule)
		strBool := "Equals(\"" + test.str + strings.Join(rulesAsString, ", ") + "\") should be " + strconv.FormatBool(test.resultBool)
		resultRule, resultBool := s.MatcherUtils.MatchAny(test.str)
		assert.Equal(s.T(), test.resultRule, resultRule, strRule)
		assert.Equal(s.T(), test.resultBool, resultBool, strBool)
	}
}


func TestMatcherUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(MatcherUtilsTestSuite))
}
