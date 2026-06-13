package utils

import (
	"strings"
	"time"
	"unicode/utf8"
)

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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const dateLayout = "2006-01-02"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DateString formats t as a calendar date in local time
func DateString(t time.Time) string {
	return t.In(time.Local).Format(dateLayout)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ParseDateString parses a calendar date in local time at midnight
func ParseDateString(value string) (time.Time, error) {
	return time.ParseInLocation(dateLayout, value, time.Local)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NextMidnight returns the next local midnight strictly after t
func NextMidnight(t time.Time) time.Time {
	local := t.In(time.Local)

	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// PreviousDateString returns the calendar date before t in local time
func PreviousDateString(t time.Time) string {
	return DateString(t.In(time.Local).AddDate(0, 0, -1))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NormalizeGroupName trims and validates a group name against maxLength runes
func NormalizeGroupName(name string, maxLength int) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrGroupName
	}

	if utf8.RuneCountInString(name) > maxLength {
		return "", ErrGroupNameTooLong
	}

	return name, nil
}
