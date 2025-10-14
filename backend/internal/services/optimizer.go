package services

import (
	"regexp"
	"strings"
)

type QueryOptimizer struct {
	prefixPattern *regexp.Regexp
}

func NewQueryOptimizer() *QueryOptimizer {
	return &QueryOptimizer{
		prefixPattern: regexp.MustCompile(`^\*+\s*`),
	}
}

func (o *QueryOptimizer) OptimizeQuery(query, searchType string) string {
	optimized := o.prefixPattern.ReplaceAllString(query, "")
	optimized = strings.TrimSpace(optimized)
	optimized = o.cleanSpaces(optimized)
	return optimized
}

func (o *QueryOptimizer) cleanSpaces(query string) string {
	spacePattern := regexp.MustCompile(`\s+`)
	result := spacePattern.ReplaceAllString(query, " ")
	return strings.TrimSpace(result)
}

func (o *QueryOptimizer) ValidateQuery(query string) (bool, string) {
	query = strings.TrimSpace(query)

	if len(query) < 2 {
		return false, "Query too short"
	}

	if len(query) > 200 {
		return false, "Query too long"
	}

	alphanumericPattern := regexp.MustCompile(`[a-zA-Z0-9]`)
	if !alphanumericPattern.MatchString(query) {
		return false, "Query must contain letters or numbers"
	}

	return true, ""
}

func (o *QueryOptimizer) IsProductQuery(query string) bool {
	return len(strings.TrimSpace(query)) >= 2
}
