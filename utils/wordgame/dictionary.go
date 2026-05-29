package wordgame

import (
	"bufio"
	"embed"
	"strings"
	"sync"
)

//go:embed dictionary
var embedFS embed.FS

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Dictionary holds the set of five-letter words accepted for picks and guesses
type Dictionary struct {
	words map[string]struct{}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var (
	global     *Dictionary
	globalOnce sync.Once
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// LoadDictionary returns the shared word list loaded once at startup
func LoadDictionary() (*Dictionary, error) {
	var err error
	globalOnce.Do(func() {
		global, err = load()
	})

	return global, err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsValidGuess reports whether guessers may submit the word
func (d *Dictionary) IsValidGuess(word string) bool {
	return d.hasWord(word)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsValidAnswer reports whether the picker may choose the word as the round answer
func (d *Dictionary) IsValidAnswer(word string) bool {
	return d.hasWord(word)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// load reads words.txt from the embedded dictionary directory
func load() (*Dictionary, error) {
	words, err := readWords("dictionary/words.txt")
	if err != nil {
		return nil, err
	}

	return &Dictionary{words: toSet(words)}, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// hasWord reports whether the word is in the dictionary
func (d *Dictionary) hasWord(word string) bool {
	_, ok := d.words[strings.ToUpper(word)]

	return ok
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// readWords loads five-letter words from a newline-delimited file
func readWords(path string) ([]string, error) {
	f, err := embedFS.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var words []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		w := strings.ToUpper(strings.TrimSpace(sc.Text()))
		if len(w) == 5 {
			words = append(words, w)
		}
	}

	return words, sc.Err()
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// toSet builds a lookup map from a word slice
func toSet(words []string) map[string]struct{} {
	m := make(map[string]struct{}, len(words))
	for _, w := range words {
		m[w] = struct{}{}
	}

	return m
}
