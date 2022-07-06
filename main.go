package main

import (
	"fmt"
	"github.com/romainvacheret/siglog/src/matching"
)

func main() {
	defaultRegistry, _ := matching.NewDefaultRuleRegistry() 
	registry := matching.NewRuleRegistry()
	matcher := matching.NewOrderedMatcher(*registry, []matching.RuleName{
		matching.RuleName("endsWithe"),
		matching.RuleName("startsWithW"),
	})
	lines := []string{"Welcom", "World", "here", "foo"}

	registry.RegisterWithRegistry("startsWithW", "startsWith", "W", defaultRegistry)
	registry.RegisterWithRegistryBefore("endsWithe", "endsWith", "e", defaultRegistry)

	fmt.Println(matcher.CountMatches(lines))
}
