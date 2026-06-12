package utils

import "strings"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Map runs a function over a slice of type T, returning a new slice of type V
func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))

	for i, t := range ts {
		result[i] = fn(t)
	}

	return result
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// FilterMap runs fn over each item in ts, collecting values where fn returns ok true
func FilterMap[T, V any](ts []T, fn func(T) (V, bool)) []V {
	result := make([]V, 0, len(ts))

	for _, t := range ts {
		if v, ok := fn(t); ok {
			result = append(result, v)
		}
	}

	return result
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// StringSplit splits a string into a slice of strings, trimming each string and removing
// empty strings
func StringSplit(s string, sep string) []string {
	var out []string

	for _, p := range strings.Split(s, sep) {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ToSet builds a string lookup set from a slice of unique keys
func ToSet(words []string) map[string]struct{} {
	m := make(map[string]struct{}, len(words))
	for _, w := range words {
		m[w] = struct{}{}
	}

	return m
}
