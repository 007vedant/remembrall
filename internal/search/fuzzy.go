package search

import (
	"remembrall/pkg/models"
	"sort"
	"strings"
)

// MatchResult represents a fuzzy search match with a score
type MatchResult struct {
	Entry *models.PasswordEntry
	Score int
}

// FuzzySearch performs fuzzy search on password entries
func FuzzySearch(entries []*models.PasswordEntry, query string) []*MatchResult {
	if query == "" {
		return nil
	}

	query = strings.ToLower(query)
	var results []*MatchResult

	for _, entry := range entries {
		score := calculateMatchScore(strings.ToLower(entry.AppName), query)
		if score > 0 {
			results = append(results, &MatchResult{
				Entry: entry,
				Score: score,
			})
		}
	}

	// Sort by score (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// calculateMatchScore returns a score for how well the app name matches the query
// Higher score = better match
func calculateMatchScore(appName, query string) int {
	if appName == query {
		return 100 // Exact match
	}

	if strings.HasPrefix(appName, query) {
		return 90 // Prefix match
	}

	if strings.Contains(appName, query) {
		return 80 // Contains match
	}

	// Check for subsequence match (e.g., "gh" matches "github")
	if isSubsequence(query, appName) {
		return 70
	}

	// Check for word boundary matches
	words := strings.Fields(strings.ReplaceAll(appName, "-", " "))
	for _, word := range words {
		if strings.HasPrefix(word, query) {
			return 60
		}
	}

	// Levenshtein distance for typos
	if len(query) >= 3 && levenshteinDistance(appName, query) <= 2 {
		return 50
	}

	return 0 // No match
}

// isSubsequence checks if query is a subsequence of target
func isSubsequence(query, target string) bool {
	queryIdx := 0
	for i := 0; i < len(target) && queryIdx < len(query); i++ {
		if target[i] == query[queryIdx] {
			queryIdx++
		}
	}
	return queryIdx == len(query)
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	// Initialize first row and column
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

// FindBestMatch returns the best match for a query, or nil if no good match
func FindBestMatch(entries []*models.PasswordEntry, query string) *models.PasswordEntry {
	results := FuzzySearch(entries, query)
	if len(results) == 0 || results[0].Score < 50 {
		return nil
	}
	return results[0].Entry
}