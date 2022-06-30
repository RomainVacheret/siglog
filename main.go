package main

import (
	"fmt"
	"github.com/romainvacheret/siglog/src/matching"
)

func main() {
	defaultRegistry, _ := matching.NewDefaultRuleRegistry() 
	registry := matching.NewRuleRegistry()
	matcher := matching.NewUnorderedMatcher(*registry)
	lines := []string{"Welcome", "World", "here", "foo"}

	registry.RegisterWithRegistry("startsWithW", "startsWith", "W", defaultRegistry)
	registry.RegisterWithRegistryBefore("endsWithe", "endsWith", "e", defaultRegistry)

	fmt.Println(matcher.CountMatches(lines))
}
