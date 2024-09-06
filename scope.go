package scopes

import (
	"strings"
)

// strech converts a slice to bigger length by adding toAdd to it
func strech(scopeSlice []string, toAdd string, toLen int) []string {
	lenDiff := toLen - len(scopeSlice)
	for i := 0; i < lenDiff; i++ {
		scopeSlice = append(scopeSlice, toAdd)
	}
	return scopeSlice
}

// MatchScopes matches two scopes using Wildcard Scope Matching Strategy (asymetric)
func MatchScopes(scopeA, scopeB string) bool {
	scopeASplit := strings.Split(scopeA, ":")
	scopeBSplit := strings.Split(scopeB, ":")
	return matchParsedScopes(scopeASplit, scopeBSplit)
}

func matchParsedScopes(scopeASplit, scopeBSplit []string) bool {
	scopeALen := len(scopeASplit)
	scopeBLen := len(scopeBSplit)

	// If scopeBLen is smaller than scopeALen and last char of scopeB is not * return false
	if scopeBLen < scopeALen && scopeBSplit[scopeBLen-1] != "*" {
		return false
		// If scopeBLen is smaller than scopeALen and last char of scopeB is * stretch scopeB To Len Of ScopeA By Adding "*"
	} else if scopeBLen < scopeALen && scopeBSplit[scopeBLen-1] == "*" {
		scopeBSplit = strech(scopeBSplit, "*", scopeALen)
		// If scopeBLen is greater than scopeALen and last char of scopeA is not * return false
	} else if scopeBLen > scopeALen && scopeASplit[scopeALen-1] != "*" {
		return false
	}

	for i := 0; i < scopeALen; i++ {
		if !(scopeASplit[i] == scopeBSplit[i] || scopeBSplit[i] == "*") {
			return false
		}
	}

	return true
}

// ScopeInAllowed is used to check if scope is allowed based on allowed scopes list
func ScopeInAllowed(scope string, allowedScopes []string) bool {
	scopeASplit := strings.Split(scope, ":")
	for _, allowedScope := range allowedScopes {
		scopeBSplit := strings.Split(allowedScope, ":")
		if matchParsedScopes(scopeASplit, scopeBSplit) {
			return true
		}
	}
	return false
}

func ScopesInAllowed(scopes []string, allowedScopes []string) bool {
	parsedAllowedScopes := ParseScopeRule(allowedScopes)
	for _, scopeRequired := range scopes {
		scopeASplit := strings.Split(scopeRequired, ":")
		for _, scopeBSplit := range parsedAllowedScopes {
			if matchParsedScopes(scopeASplit, scopeBSplit) {
				return true
			}
		}
	}
	return false
}

func scopesInParsedAllowed(scopes []string, parsedAllowedScopes [][]string) bool {
	for _, scopeRequired := range scopes {
		scopeASplit := strings.Split(scopeRequired, ":")
		for _, scopeBSplit := range parsedAllowedScopes {
			if matchParsedScopes(scopeASplit, scopeBSplit) {
				return true
			}
		}
	}
	return false
}

func ParseScopeRule(allowedScopes []string) [][]string {
	parsedAllowedScopes := make([][]string, len(allowedScopes))
	for index, allowedScope := range allowedScopes {
		parsedAllowedScopes[index] = strings.Split(allowedScope, ":")
	}
	return parsedAllowedScopes
}
