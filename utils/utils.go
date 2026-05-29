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
