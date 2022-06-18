package main

import (
	"fmt"
	"siglog/src/matching"
)

func main() {
	registry := matching.NewRuleRegistry()
	matcher := matching.NewMatcher(*registry)

	registry.RegisterWith("startsWithW", "startsWith", "W")
	registry.RegisterWithBefore("endsWithe", "endsWith", "e")

	fmt.Println(matcher.MatchString("startsWithW", "Welcome"))
	fmt.Println(matcher.MatchString("endsWithe", "Welcome"))
}